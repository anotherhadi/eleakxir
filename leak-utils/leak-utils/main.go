package main

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/misc"
	"github.com/anotherhadi/eleakxir/leak-utils/parquet"
	"github.com/anotherhadi/eleakxir/leak-utils/settings"
	"github.com/charmbracelet/log"
	_ "github.com/marcboeker/go-duckdb/v2"
	flag "github.com/spf13/pflag"
)

func main() {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal("Failed to open DuckDB", "error", err)
	}
	defer db.Close()
	lu := settings.LeakUtils{
		Db: db,
	}
	actions := []string{
		"cleanParquet",
		"infoParquet",
		// Csv
		"csvToParquet",
		// Misc
		"mergeFiles",
		"removeUrlSchemeFromUlp",
	}

	if len(os.Args) < 2 {
		fmt.Println(settings.Muted.Render("Usage: "), settings.Accent.Render(os.Args[0], "<action>"))
		fmt.Println(settings.Muted.Render("Actions: "), settings.Base.Render(strings.Join(actions, ", ")))
		return
	}
	action := os.Args[1]
	if !slices.Contains(actions, action) {
		log.Fatal("Unknown action", "action", action)
	}

	switch action {
	case "cleanParquet":
		var inputFile *string = flag.StringP("input", "i", "", "Input Parquet file")
		var outputFile *string = flag.StringP("output", "o", "", "Output Parquet file")
		var compression *string = flag.StringP("compression", "c", "ZSTD", "Compression codec (UNCOMPRESSED, SNAPPY, GZIP, BROTLI, LZ4, ZSTD)")
		var skipLineFormating *bool = flag.BoolP("skip-line-formating", "s", false, "Skip line formating")
		var deleteFirstRow *bool = flag.Bool("delete-first-row", false, "Delete first row")
		var debug *bool = flag.Bool("debug", false, "Debug mode")
		var noColors *bool = flag.Bool("no-colors", false, "Remove all colors")
		var printQuery *bool = flag.BoolP("print-query", "p", false, "Print the query instead of executing it")
		flag.Parse()
		if *inputFile == "" || *outputFile == "" {
			log.Fatal("Input and output files are required")
		}
		if *noColors {
			settings.DisableColors()
		}
		lu.Compression = *compression
		lu.Debug = *debug
		err := parquet.CleanParquet(lu, *inputFile, *outputFile, *skipLineFormating, *deleteFirstRow, *printQuery)
		if err != nil {
			log.Fatal("Failed to clean Parquet file", "error", err)
		}
		return
	case "infoParquet":
		var inputFile *string = flag.StringP("input", "i", "", "Input Parquet file")
		var debug *bool = flag.Bool("debug", false, "Debug mode")
		var noColors *bool = flag.Bool("no-colors", false, "Remove all colors")
		flag.Parse()
		if *inputFile == "" {
			log.Fatal("Input files are required")
		}
		if *noColors {
			settings.DisableColors()
		}
		lu.Debug = *debug
		err := parquet.InfoParquet(lu, *inputFile)
		if err != nil {
			log.Fatal("Failed to read Parquet file", "error", err)
		}
		return
	case "csvToParquet":
		var inputFile *string = flag.StringP("input", "i", "", "Input Parquet file")
		var outputFile *string = flag.StringP("output", "o", "", "Output Parquet file")
		var strict *bool = flag.Bool("strict", true, "Strict mode for Duckdb")
		var compression *string = flag.StringP("compression", "c", "ZSTD", "Compression codec (UNCOMPRESSED, SNAPPY, GZIP, BROTLI, LZ4, ZSTD)")
		var noColors *bool = flag.Bool("no-colors", false, "Remove all colors")
		var debug *bool = flag.Bool("debug", false, "Debug mode")
		flag.Parse()
		if *inputFile == "" || *outputFile == "" {
			log.Fatal("Input and output files are required")
		}
		if *noColors {
			settings.DisableColors()
		}
		lu.Compression = *compression
		lu.Debug = *debug
		err := misc.CsvToParquet(lu, *inputFile, *outputFile, *strict)
		if err != nil {
			log.Fatal("Failed to transform Csv file", "error", err)
		}
		return
	case "mergeFiles":
		var inputFiles *[]string = flag.StringArrayP("inputs", "i", []string{}, "Input Parquet files")
		var outputFile *string = flag.StringP("output", "o", "", "Output Parquet file")
		var noColors *bool = flag.Bool("no-colors", false, "Remove all colors")
		var debug *bool = flag.Bool("debug", false, "Debug mode")
		flag.Parse()
		if len(*inputFiles) == 0 || *outputFile == "" {
			log.Fatal("Inputs and output files are required")
		}
		if *noColors {
			settings.DisableColors()
		}
		lu.Debug = *debug
		err := misc.MergeFiles(lu, *outputFile, *inputFiles...)
		if err != nil {
			log.Fatal("Failed to merge files", "error", err)
		}
		return
	case "removeUrlSchemeFromUlp":
		var inputFile *string = flag.StringP("input", "i", "", "Input Parquet file")
		var noColors *bool = flag.Bool("no-colors", false, "Remove all colors")
		var debug *bool = flag.Bool("debug", false, "Debug mode")
		flag.Parse()
		if *inputFile == "" {
			log.Fatal("Input files are required")
		}
		if *noColors {
			settings.DisableColors()
		}
		lu.Debug = *debug
		err := misc.RemoveUrlSchemeFromUlp(lu, *inputFile)
		if err != nil {
			log.Fatal("Failed to remove ULP Url schemes", "error", err)
		}
		return
	}
}
