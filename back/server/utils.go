package server

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getParquetColumns(s Server, path string) []string {
	query := fmt.Sprintf("DESCRIBE SELECT * FROM '%s';", path)

	rows, err := s.Duckdb.Query(query)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	var columns []string

	for rows.Next() {
		var columnName string
		var columnType string
		var nullable string

		var key sql.NullString
		var defaultValue sql.NullString
		var extra sql.NullString

		if err := rows.Scan(&columnName, &columnType, &nullable, &key, &defaultValue, &extra); err != nil {
			return []string{}
		}
		columns = append(columns, columnName)
	}

	if err = rows.Err(); err != nil {
		return []string{}
	}

	if len(columns) == 0 {
		return []string{}
	}

	return columns
}

func getParquetLength(s Server, path string) uint64 {
	query := fmt.Sprintf("SELECT COUNT(*) FROM '%s';", path)

	row := s.Duckdb.QueryRow(query)

	var count uint64
	if err := row.Scan(&count); err != nil {
		return 0
	}

	return count
}

// Walk through the given folder and its subfolders to find all parquet files
// Return a list of path
func getAllParquetFiles(folders []string) []string {
	var paths []string
	for _, baseDir := range folders {
		_ = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".parquet") {
				return err
			}
			paths = append(paths, path)
			return nil
		})
	}
	return paths
}

func getFileSize(path string) uint64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return uint64(info.Size() / (1024 * 1024)) // MB
}

func FormatParquetName(path string) string {
	_, file := filepath.Split(path)
	fileName := strings.TrimSuffix(file, ".parquet")

	parts := strings.Split(fileName, "-")
	sourceName := parts[0]
	var blocks []string

	for _, part := range parts[1:] {
		if strings.HasPrefix(part, "date_") {
			dateStr := strings.TrimPrefix(part, "date_")
			dateStr = strings.ReplaceAll(dateStr, "_", "/")
			blocks = append(blocks, fmt.Sprintf("date: %s", dateStr))
		} else if strings.HasPrefix(part, "source_") {
			sourceStr := strings.TrimPrefix(part, "source_")
			blocks = append(blocks, fmt.Sprintf("source: %s", sourceStr))
		} else if strings.HasPrefix(part, "notes_") {
			noteStr := strings.TrimPrefix(part, "notes_")
			noteStr = strings.ReplaceAll(noteStr, "_", " ")
			blocks = append(blocks, noteStr)
		}
	}

	sourceName = strings.ReplaceAll(sourceName, "_", " ")
	sourceWords := strings.Fields(sourceName)
	for i, word := range sourceWords {
		if len(word) > 0 {
			sourceWords[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	formattedSourceName := strings.Join(sourceWords, " ")

	if len(blocks) > 0 {
		return fmt.Sprintf("%s (%s)", formattedSourceName, strings.Join(blocks, ", "))
	}

	return formattedSourceName
}

func createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}
