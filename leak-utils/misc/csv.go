package misc

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
	"github.com/charmbracelet/log"
)

func CsvToParquet(lu settings.LeakUtils, inputFile string, outputFile string, strict bool) error {
	delimiter := getDelimiter(inputFile)

	hasHeader, err := csvHasHeader(inputFile, delimiter)
	if err != nil {
		return err
	}
	header := "true"
	if !hasHeader {
		header = "false"
	}
	strictMode := "true"
	if !strict {
		strictMode = "false"
	}

	query := fmt.Sprintf(`CREATE TABLE my_table AS FROM read_csv_auto('%s', HEADER=%s, delim='%s', ignore_errors=true, all_varchar=true, null_padding=true, strict_mode=%s);
		COPY my_table TO '%s' (FORMAT 'parquet', COMPRESSION '%s', ROW_GROUP_SIZE 200_000);`,
		inputFile, header, delimiter, strictMode, outputFile, lu.Compression)

	if lu.Debug {
		log.Info("Detected delimiter", "delimiter", delimiter)
		log.Info("CSV header detection", "hasHeader", hasHeader)
		log.Info("Executing query", "query", query)
	}

	_, err = lu.Db.Exec(query)

	if lu.Debug {
		log.Info("Finished executing query")
	}

	return err
}

func getDelimiter(inputFile string) string {
	lines, err := getNLine(inputFile, 10, 0)
	if err != nil {
		log.Warn("Failed to read CSV file to determine delimiter, defaulting to comma", "error", err)
		return ","
	}

	delimiterCounts := map[string]int{
		",":  0,
		";":  0,
		"\t": 0,
		"|":  0,
		":":  0,
	}

	for _, line := range lines {
		for d := range delimiterCounts {
			delimiterCounts[d] += strings.Count(line, d)
		}
	}

	maxCount := 0
	delimiter := ","

	for d, count := range delimiterCounts {
		if count > maxCount {
			maxCount = count
			delimiter = d
		}
	}

	return delimiter
}

func csvHasHeader(inputFile, delimiter string) (hasHeader bool, err error) {
	firstRow, err := getFirstRowCsv(inputFile, delimiter)
	if err != nil {
		return false, err
	}
	for i, col := range firstRow {
		col = strings.ReplaceAll(col, "\"", "")
		col = strings.ReplaceAll(col, " ", "")
		col = strings.ReplaceAll(col, "-", "")
		col = strings.ReplaceAll(col, "_", "")
		col = strings.ReplaceAll(col, ".", "")
		firstRow[i] = strings.ToLower(strings.TrimSpace(col))
	}
	knownHeaders := []string{"email", "password", "username", "phone", "lastname", "firstname", "mail", "addresse", "nom", "id"}
	for _, knownHeader := range knownHeaders {
		if slices.Contains(firstRow, knownHeader) {
			return true, nil
		}
	}
	return false, nil
}

func getNLine(inputFile string, n, offset int) (lines []string, err error) {
	if n <= 0 {
		return nil, nil
	}

	if offset < 0 {
		offset = 0
	}

	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine <= offset {
			continue
		}

		lines = append(lines, scanner.Text())
		if len(lines) >= n {
			break
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return nil, err
	}

	return lines, nil
}

func getFirstRowCsv(inputFile, delimiter string) (row []string, err error) {
	rows, err := getFirstNRowsCsv(inputFile, 1, delimiter)
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found in CSV")
	}
	return rows[0], err
}

func getFirstNRowsCsv(inputFile string, n int, delimiter string) (rows [][]string, err error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	if len(delimiter) != 1 {
		return nil, fmt.Errorf("delimiter must be a single character, got %q", delimiter)
	}
	reader.Comma = rune(delimiter[0])

	for i := 0; i < n; i++ {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to read CSV: %w", err)
		}
		rows = append(rows, row)
	}

	return rows, nil
}
