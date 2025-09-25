package dataleak

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/anotherhadi/eleakxir/backend/server"
	"github.com/charmbracelet/log"
)

type LeakResult struct {
	Duration time.Duration
	Rows     []map[string]string
	Error    string
	LimitHit bool // Whether the search hit the limit
}

func Search(s *server.Server, queryText, column string, exactMatch bool) LeakResult {
	if len(*(s.Dataleaks)) == 0 {
		return LeakResult{
			Error: "No dataleak configured",
		}
	}
	now := time.Now()
	result := LeakResult{}

	sqlQuery := buildSqlQuery(s, queryText, column, exactMatch)

	if s.Settings.Debug {
		log.Info("New query:", "query", sqlQuery)
	}
	rows, err := s.Duckdb.Query(sqlQuery)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		result.Error = err.Error()
		return result
	}

	rawResult := make([][]byte, len(cols))
	dest := make([]any, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			result.Error = err.Error()
			return result
		}

		rowMap := make(map[string]string)
		for i, colName := range cols {
			if rawResult[i] == nil || colName == "" {
				continue
			}
			if colName == "filename" {
				rowMap["source"] = server.FormatParquetName(string(rawResult[i]))
				continue
			}
			rowMap[colName] = string(rawResult[i])
		}
		result.Rows = append(result.Rows, rowMap)
	}

	if err = rows.Err(); err != nil {
		result.Error = err.Error()
		return result
	}

	if len(result.Rows) >= s.Settings.Limit {
		result.LimitHit = true
	}

	result.Rows = removeDuplicateMaps(result.Rows)

	result.Duration = time.Since(now)
	return result
}

func removeDuplicateMaps(maps []map[string]string) []map[string]string {
	seen := make(map[string]struct{})
	result := []map[string]string{}

	for _, m := range maps {
		// Create a unique key for the map by concatenating its key-value pairs
		var sb strings.Builder
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		slices.Sort(keys) // Sort keys to ensure consistent order
		for _, k := range keys {
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(m[k])
			sb.WriteString(";")
		}
		key := sb.String()

		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, m)
		}
	}

	return result
}

func buildSqlQuery(s *server.Server, queryText, column string, exactMatch bool) string {
	// Normalize "name" -> "full_name"
	if strings.EqualFold(column, "name") {
		column = "full_name"
	}

	// Step 1: Determine candidate columns to search
	var candidateColumns []string
	if column == "all" || column == "" {
		// Use base columns if "all" or empty
		candidateColumns = s.Settings.BaseColumns
	} else {
		// Otherwise, only search the given column
		candidateColumns = []string{column}
	}

	// Step 2: Collect all available columns across dataleaks
	allColumns := make([]string, 0)
	seen := make(map[string]struct{})
	for _, dataleak := range *s.Dataleaks {
		for _, col := range dataleak.Columns {
			if _, ok := seen[col]; !ok {
				seen[col] = struct{}{}
				allColumns = append(allColumns, col)
			}
		}
	}

	// Step 3: Resolve which columns should actually be used in the WHERE clause
	var columnsFiltered []string
	if strings.EqualFold(column, "full_text") {
		// "full_text" means search across all columns
		columnsFiltered = allColumns
	} else {
		for _, candidate := range candidateColumns {
			for _, available := range allColumns {
				// Exact match (case-insensitive)
				if strings.EqualFold(available, candidate) {
					columnsFiltered = append(columnsFiltered, available)
					continue
				}
				// Match columns ending with "_<candidate>"
				if strings.HasSuffix(strings.ToLower(available), "_"+strings.ToLower(candidate)) {
					columnsFiltered = append(columnsFiltered, available)
				}
			}
		}
	}

	limit := strconv.Itoa(s.Settings.Limit)
	from := getFromClause(s)

	if len(columnsFiltered) == 0 {
		return fmt.Sprintf("SELECT * FROM %s LIMIT %s", from, limit)
	}

	where := getWhereClause(queryText, columnsFiltered, exactMatch)

	return fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT %s", from, where, limit)
}

func getWhereClause(queryText string, columns []string, exactMatch bool) string {
	terms := strings.Fields(queryText)
	var andClauses []string

	for _, term := range terms {
		var orClausesForTerm []string
		termEscaped := strings.ReplaceAll(term, "'", "''")

		startsWith := false
		endsWith := false
		if strings.HasPrefix(termEscaped, "^") {
			startsWith = true
			termEscaped = strings.TrimPrefix(termEscaped, "^")
		}
		if strings.HasSuffix(termEscaped, "$") {
			endsWith = true
			termEscaped = strings.TrimSuffix(termEscaped, "$")
		}

		termEscapedILike := strings.ReplaceAll(termEscaped, "_", "\\_")
		termEscapedILike = strings.ReplaceAll(termEscapedILike, "%", "\\%")
		for _, col := range columns {
			if exactMatch || (startsWith && endsWith) {
				orClausesForTerm = append(orClausesForTerm, fmt.Sprintf("\"%s\" ILIKE '%s' ESCAPE '\\'", col, strings.ToLower(termEscapedILike)))
			} else if startsWith {
				orClausesForTerm = append(orClausesForTerm, fmt.Sprintf("\"%s\" ILIKE '%s%%' ESCAPE '\\'", col, strings.ToLower(termEscapedILike)))
			} else if endsWith {
				orClausesForTerm = append(orClausesForTerm, fmt.Sprintf("\"%s\" ILIKE '%%%s' ESCAPE '\\'", col, strings.ToLower(termEscapedILike)))
			} else {
				orClausesForTerm = append(orClausesForTerm, fmt.Sprintf("\"%s\" ILIKE '%%%s%%' ESCAPE '\\'", col, strings.ToLower(termEscapedILike)))
			}
		}
		andClauses = append(andClauses, "("+strings.Join(orClausesForTerm, " OR ")+")")
	}
	return strings.Join(andClauses, " AND ")
}

func getFromClause(s *server.Server) string {
	parquets := []string{}
	for _, dataleak := range *s.Dataleaks {
		parquets = append(parquets, "'"+dataleak.Path+"'")
	}
	return fmt.Sprintf("read_parquet([%s], union_by_name=true, filename=true)", strings.Join(parquets, ", "))
}
