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
	limit := strconv.Itoa(s.Settings.Limit)
	from := getFromClause(s)
	if column == "name" {
		column = "full_name"
	}
	columns := []string{column}
	if column == "all" || column == "" {
		columns = s.Settings.BaseColumns
	}
	columnsFiltered := []string{}
	allColumns := []string{}
	// TODO: Add columns that ends with _col aswell
	for _, dataleak := range *s.Dataleaks {
		for _, col := range dataleak.Columns {
			if !slices.Contains(allColumns, col) {
				allColumns = append(allColumns, col)
			}
		}
	}
	if column == "full_text" {
		columnsFiltered = allColumns
	} else {
		for _, col := range columns {
			if slices.Contains(allColumns, col) {
				columnsFiltered = append(columnsFiltered, col)
			}
		}
	}

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

		for _, col := range columns {
			if exactMatch {
				termEscapedILike := strings.ReplaceAll(termEscaped, "_", "\\_")
				termEscapedILike = strings.ReplaceAll(termEscapedILike, "%", "\\%")
				orClausesForTerm = append(orClausesForTerm, fmt.Sprintf("\"%s\" ILIKE '%s' ESCAPE '\\'", col, strings.ToLower(termEscapedILike)))
			} else {
				// Escape characters for ILIKE
				termEscapedILike := strings.ReplaceAll(termEscaped, "_", "\\_")
				termEscapedILike = strings.ReplaceAll(termEscapedILike, "%", "\\%")
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

func castAllColumns(cols []string) []string {
	casted := make([]string, len(cols))
	for i, col := range cols {
		casted[i] = fmt.Sprintf("cast(%s as text)", col)
	}
	return casted
}
