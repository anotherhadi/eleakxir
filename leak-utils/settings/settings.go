package settings

import "database/sql"

type LeakUtils struct {
	Debug       bool
	Compression string // Compression of the output file (e.g., "SNAPPY", "ZSTD", "NONE" or "")
	Db          *sql.DB
}
