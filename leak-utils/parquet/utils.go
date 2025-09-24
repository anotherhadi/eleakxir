package parquet

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// GetColumns retrieves the column names from the Parquet file.
func GetColumns(db *sql.DB, filepath string) ([]string, error) {
	// Create a view from the parquet file
	query := fmt.Sprintf("CREATE OR REPLACE VIEW parquet_view AS SELECT * FROM read_parquet('%s')", filepath)
	_, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create view: %w", err)
	}

	// Get column information
	rows, err := db.Query("DESCRIBE parquet_view")
	if err != nil {
		return nil, fmt.Errorf("failed to describe view: %w", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var colName, colType, nullable, key, defaultVal, extra sql.NullString
		err := rows.Scan(&colName, &colType, &nullable, &key, &defaultVal, &extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if colName.Valid {
			columns = append(columns, colName.String)
		}
	}

	return columns, nil
}

// getFirstNRows retrieves the first N rows from the Parquet file.
func getFirstNRows(db *sql.DB, inputFile string, n int) ([][]string, error) {
	query := fmt.Sprintf("SELECT * FROM read_parquet('%s') LIMIT %d", inputFile, n)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query parquet file: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results [][]string
	for rows.Next() {
		values := make([]sql.NullString, len(cols))
		valuePtrs := make([]any, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var row []string
		for _, val := range values {
			if val.Valid {
				row = append(row, val.String)
			} else {
				row = append(row, "NULL")
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// countRows counts the number of rows in the Parquet file.
func countRows(db *sql.DB, inputFile string) (int64, error) {
	var count int64
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM read_parquet('%s')", inputFile)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows: %w", err)
	}
	return count, nil
}

// formatWithSpaces formats an integer with spaces as thousand separators.
func formatWithSpaces(n int64) string {
	s := strconv.FormatInt(n, 10)

	var b strings.Builder
	l := len(s)
	for i, c := range s {
		if i != 0 && (l-i)%3 == 0 {
			b.WriteRune(' ')
		}
		b.WriteRune(c)
	}
	return b.String()
}
