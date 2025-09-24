package misc

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/parquet"
	"github.com/anotherhadi/eleakxir/leak-utils/settings"
)

// Count the line with "@" in a file
func CountLinesWithAt(lu settings.LeakUtils, inputFile string) (nAt, nLines int, err error) {
	if strings.HasSuffix(inputFile, ".parquet") {
		return countRowsWithAtInParquet(lu, inputFile)
	}

	in, err := os.Open(inputFile)
	if err != nil {
		return 0, 0, err
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	countAt := 0
	countLine := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "@") {
			countAt++
		}
		countLine++
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return countAt, countLine, nil
}

func countRowsWithAtInParquet(lu settings.LeakUtils, inputFile string) (nAt, nLine int, err error) {
	cols, err := parquet.GetColumns(lu.Db, inputFile)
	if err != nil {
		return 0, 0, err
	}
	if len(cols) == 0 {
		return 0, 0, nil
	}

	whereParts := []string{}
	for _, col := range cols {
		whereParts = append(whereParts, fmt.Sprintf("%s LIKE '%%@%%'", col))
	}
	whereClause := strings.Join(whereParts, " OR ")

	query := fmt.Sprintf("SELECT COUNT(*) FROM read_parquet('%s') WHERE %s", inputFile, whereClause)
	var countAt int
	err = lu.Db.QueryRow(query).Scan(&countAt)
	if err != nil {
		return 0, 0, err
	}

	query = fmt.Sprintf("SELECT COUNT(*) FROM read_parquet('%s')", inputFile)
	var countLine int
	err = lu.Db.QueryRow(query).Scan(&countLine)
	if err != nil {
		return 0, 0, err
	}

	return countAt, countLine, nil
}
