package parquet

import (
	"fmt"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
)

// If there is no full_name but there is last_name and first_name, create full_name
// If there is no full_name, no last_name or no first_name, but there is name, rename name to full_name
func addFullname(operations []ColumnOperation) []ColumnOperation {
	hasFullName := false
	hasFirstName := false
	hasLastName := false
	hasName := false
	for _, op := range operations {
		if op.Action != "drop" {
			if op.NewName == "full_name" {
				hasFullName = true
			} else if op.NewName == "first_name" {
				hasFirstName = true
			} else if op.NewName == "last_name" {
				hasLastName = true
			} else if op.NewName == "name" {
				hasName = true
			}
		}
	}
	if hasFullName {
		return operations
	}
	if hasFirstName && hasLastName {
		operations = append(operations, ColumnOperation{
			OriginalName: "first_name || ' ' || last_name",
			NewName:      "full_name",
			Action:       "rename",
		})
		fmt.Println(settings.Muted.Render("\nAdding new column 'full_name' as concatenation of 'first_name' and 'last_name'."))
		return operations
	}
	if hasName {
		for i, op := range operations {
			if op.NewName == "name" && op.Action != "drop" {
				operations[i].NewName = "full_name"
				fmt.Println(settings.Muted.Render("\nRenaming column 'name' to 'full_name'."))
				return operations
			}
		}
	}
	if hasFirstName {
		operations = append(operations, ColumnOperation{
			OriginalName: "first_name",
			NewName:      "full_name",
			Action:       "rename",
		})
		fmt.Println(settings.Muted.Render("\nAdding new column 'full_name' from 'first_name'."))
		return operations
	}
	if hasLastName {
		operations = append(operations, ColumnOperation{
			OriginalName: "last_name",
			NewName:      "full_name",
			Action:       "rename",
		})
		fmt.Println(settings.Muted.Render("\nAdding new column 'full_name' from 'last_name'."))
		return operations
	}

	return operations
}

// formatColumnName formats a column name to be SQL-compliant.
func formatColumnName(columnName string) string {
	columnName = strings.TrimSpace(columnName)
	columnName = strings.ToLower(columnName)
	columnName = strings.Join(strings.Fields(columnName), "_")
	columnName = strings.ReplaceAll(columnName, "\"", "")
	columnName = strings.ReplaceAll(columnName, "'", "")
	columnName = strings.ReplaceAll(columnName, " ", "_")
	columnName = strings.ReplaceAll(columnName, "-", "_")
	// Only keep a-z, 0-9 and _
	var formatted strings.Builder
	for _, r := range columnName {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			formatted.WriteRune(r)
		}
	}
	columnName = formatted.String()
	columnName = strings.TrimPrefix(columnName, "_")
	columnName = strings.TrimSuffix(columnName, "_")
	return columnName
}

// formatColumns applies specific formatting rules to column operations.
func formatColumns(operations []ColumnOperation) []ColumnOperation {
	formatedOperations := []ColumnOperation{}
	for _, op := range operations {
		if op.NewName == "phone" || strings.HasSuffix(op.NewName, "_phone") {
			op.OriginalName = "REGEXP_REPLACE(" + op.OriginalName + ", '[^0-9]', '')"
		} else if op.NewName == "email" || strings.HasSuffix(op.NewName, "_email") {
			op.OriginalName = "REGEXP_REPLACE(LOWER(TRIM(" + op.OriginalName + ")), '[^a-z0-9._@-]', '')"
		}
		formatedOperations = append(formatedOperations, op)
	}
	return formatedOperations
}
