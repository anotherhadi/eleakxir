package parquet

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
	"github.com/charmbracelet/log"
)

type Parquet struct {
	Filepath    string
	Filename    string
	Columns     []string
	Sample      [][]string
	NRows       int64
	Compression string // Compression of the output file (e.g., "SNAPPY", "ZSTD", "NONE" or "")
}

type ColumnOperation struct {
	OriginalName string
	NewName      string
	Action       string // "keep", "rename", "drop"
}

func (parquet Parquet) PrintParquet() {
	fmt.Println(settings.Header.Render(parquet.Filename) + "\n")
	fmt.Println(settings.Accent.Render("File path:"), settings.Base.Render(parquet.Filepath))
	fmt.Println(settings.Accent.Render("Number of columns:"), settings.Base.Render(fmt.Sprintf("%d", len(parquet.Columns))))
	fmt.Println(settings.Accent.Render("Number of rows:"), settings.Base.Render(formatWithSpaces(parquet.NRows)))
	fmt.Println()
	fmt.Println(settings.Accent.Render(strings.Join(parquet.Columns, " | ")))
	for _, row := range parquet.Sample {
		fmt.Println(settings.Base.Render(strings.Join(row, " | ")))
	}
}

func InfoParquet(lu settings.LeakUtils, inputFile string) error {
	parquet, err := GetParquet(lu.Db, inputFile)
	if err != nil {
		return err
	}

	parquet.PrintParquet()
	return nil
}

func CleanParquet(lu settings.LeakUtils, inputFile, outputFile string, skipLineFormating, deleteFirstRow, printQuery bool) error {
	input, err := GetParquet(lu.Db, inputFile)
	if err != nil {
		return err
	}
	input.PrintParquet()
	columnOps := configureColumns(*input, skipLineFormating)
	output := Parquet{
		Filepath:    outputFile,
		Compression: lu.Compression,
	}
	err = transformParquet(lu, *input, output, columnOps, deleteFirstRow, printQuery)
	return err
}

func configureColumns(input Parquet, skipLineFormating bool) []ColumnOperation {
	reader := bufio.NewReader(os.Stdin)
	var operations []ColumnOperation

	fmt.Println()
	fmt.Println(settings.Base.Render("For each column, choose an action:"))
	fmt.Println(settings.Base.Render("  [k] Keep"))
	fmt.Println(settings.Base.Render("  [r] Rename"))
	fmt.Println(settings.Base.Render("  [d] Drop/Delete"))
	fmt.Println(settings.Base.Render("  [s] Suggested"))
	fmt.Println(settings.Base.Render("  [b] Go back"))
	fmt.Println()

	for i := 0; i < len(input.Columns); i++ {
		col := input.Columns[i]
		suggestion := getSuggestion(col)

		for {
			fmt.Println(settings.Muted.Render("\nColumn:"), settings.Accent.Render(col))
			if suggestion != "" {
				fmt.Println(settings.Alert.Render("Suggested action: Rename to '" + suggestion + "'"))
			}
			fmt.Print(settings.Base.Render("[k/r/d/s/b]: "))

			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading input: %v", err)
				continue
			}
			input = strings.TrimSpace(strings.ToLower(input))

			op := ColumnOperation{
				OriginalName: col,
				NewName:      col,
				Action:       "keep",
			}

			switch input {
			case "b", "back":
				if i > 0 {
					i -= 2
					if len(operations) > 0 {
						operations = operations[:len(operations)-1]
					}
					fmt.Println(settings.Muted.Render("Going back to the previous column..."))
				} else {
					fmt.Println(settings.Muted.Render("Already at the first column, cannot go back further."))
					continue
				}
				goto nextColumn

			case "r", "rename":
				fmt.Print(settings.Base.Render("Enter new name: "))
				newName, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("Error reading new name: %v", err)
					continue
				}
				newName = strings.TrimSpace(newName)
				if newName != "" {
					op.OriginalName = "\"" + op.OriginalName + "\""
					op.NewName = formatColumnName(newName)
					op.Action = "rename"
					operations = append(operations, op)
					goto nextColumn
				} else {
					fmt.Println(settings.Muted.Render("Invalid name, please try again."))
					continue
				}

			case "s", "suggested":
				if suggestion != "" {
					op.OriginalName = "\"" + op.OriginalName + "\""
					op.NewName = formatColumnName(suggestion)
					op.Action = "rename"
				} else {
					fmt.Println(settings.Muted.Render("No valid suggestion available"))
					continue
				}
				operations = append(operations, op)
				goto nextColumn

			case "d", "drop", "delete":
				op.Action = "drop"
				operations = append(operations, op)
				goto nextColumn

			case "k", "keep", "":
				op.OriginalName = "\"" + op.OriginalName + "\""
				if input == "" && suggestion != "" {
					op.NewName = formatColumnName(suggestion)
				} else {
					op.NewName = formatColumnName(op.NewName)
				}
				op.Action = "rename"
				operations = append(operations, op)
				goto nextColumn

			default:
				fmt.Println(settings.Muted.Render("Invalid choice, please enter [k/r/d/s/b]."))
				continue
			}
		}
	nextColumn:
		lastOp := operations[len(operations)-1]
		switch lastOp.Action {
		case "rename":
			if formatColumnName(lastOp.OriginalName) == lastOp.NewName {
				fmt.Printf(settings.Muted.Render("Keeping column '%s' as is.\n"), lastOp.OriginalName)
			} else {
				fmt.Printf(settings.Muted.Render("Renaming column '%s' to '%s'.\n"), lastOp.OriginalName, lastOp.NewName)
			}
		case "drop":
			fmt.Printf(settings.Muted.Render("Dropping column '%s'.\n"), lastOp.OriginalName)
		}
	}
	if !skipLineFormating {
		operations = formatColumns(operations)
	}
	operations = addFullname(operations)

	return operations
}

func transformParquet(lu settings.LeakUtils, input, output Parquet, operations []ColumnOperation, deleteFirstRow, printQuery bool) error {
	var selectClauses []string
	var columnsLength []string
	hasColumns := false

	for _, op := range operations {
		escapedOriginalName := escapeColumnName(op.OriginalName)

		if op.Action != "drop" {
			hasColumns = true

			originalSelectName := op.OriginalName
			if op.Action == "rename" {
				originalSelectName = op.OriginalName
				selectClauses = append(selectClauses, fmt.Sprintf("%s AS \"%s\"", originalSelectName, op.NewName))
			} else {
				selectClauses = append(selectClauses, originalSelectName)
			}

			columnsLength = append(columnsLength, fmt.Sprintf("COALESCE(LENGTH(\"%s\"),0)", escapedOriginalName))
		} else {
			columnsLength = append(columnsLength, fmt.Sprintf("COALESCE(LENGTH(\"%s\"),0)", escapedOriginalName))
		}
	}

	if !hasColumns {
		return fmt.Errorf("no columns selected for output")
	}

	selectClause := strings.Join(selectClauses, ", ")
	compression := ""
	if output.Compression != "" {
		compression = ", COMPRESSION '" + output.Compression + "'"
	}

	allowedRowSize := 30 * len(input.Columns)
	offset := ""
	if deleteFirstRow {
		offset = "OFFSET 1"
	}

	query := fmt.Sprintf(`
		COPY (
			SELECT %s 
			FROM read_parquet('%s')
			WHERE (%s) < %d %s
		) TO '%s' (FORMAT PARQUET, ROW_GROUP_SIZE 200_000 %s)
	`, selectClause, input.Filepath, strings.Join(columnsLength, "+"), allowedRowSize, offset, output.Filepath, compression)

	if printQuery {
		fmt.Println(settings.Base.Render("\nQuery:"))
		fmt.Println(settings.Accent.Render(strings.ReplaceAll(strings.TrimSpace(query), "\t", "")))
		return nil
	}
	if lu.Debug {
		fmt.Println(settings.Base.Render("\nQuery:"))
		fmt.Println(settings.Accent.Render(strings.ReplaceAll(strings.TrimSpace(query), "\t", "")))
	}

	fmt.Println(settings.Base.Render("\nTransforming and writing to output parquet..."))
	_, err := lu.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute transformation: %w", err)
	}
	fmt.Println(settings.Base.Render("Transformation complete!\n"))

	newParquet, err := GetParquet(lu.Db, output.Filepath)
	if err != nil {
		return err
	}
	newParquet.PrintParquet()

	return nil
}

func GetParquet(db *sql.DB, inputFile string) (parquet *Parquet, err error) {
	parquet = &Parquet{}
	parquet.Filepath = inputFile

	parquet.Columns, err = GetColumns(db, inputFile)
	if err != nil {
		return
	}
	parquet.NRows, err = countRows(db, inputFile)
	if err != nil {
		return
	}
	parquet.Sample, err = getFirstNRows(db, inputFile, 6)
	if err != nil {
		return
	}

	n := strings.LastIndex(inputFile, "/")
	if n == -1 {
		parquet.Filename = inputFile
	} else {
		parquet.Filename = inputFile[n+1:]
	}

	return
}

func escapeColumnName(name string) string {
	return strings.ReplaceAll(name, "\"", "\"\"")
}
