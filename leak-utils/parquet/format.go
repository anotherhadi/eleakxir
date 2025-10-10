package parquet

import (
	"strings"
)

// formatColumnName formats a column name to be SQL-compliant.
func formatColumnName(columnName string) string {
	columnName = strings.TrimSpace(columnName)
	columnName = strings.ToLower(columnName)
	columnName = strings.Join(strings.Fields(columnName), "_")
	columnName = strings.ReplaceAll(columnName, "\"", "")
	columnName = strings.ReplaceAll(columnName, "'", "")
	columnName = strings.ReplaceAll(columnName, " ", "_")
	columnName = strings.ReplaceAll(columnName, "-", "_")
	columnName = strings.ReplaceAll(columnName, ".", "_")
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
