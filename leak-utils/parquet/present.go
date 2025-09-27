package parquet

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
)

type dataleak struct {
	Path    string
	Name    string
	Columns []string
	Length  uint64
	Size    uint64
}

func Present(cacheFile string) error {
	dataleaks := []dataleak{}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return fmt.Errorf("error reading cache file: %w", err)
	}
	if err := json.Unmarshal(data, &dataleaks); err != nil {
		return fmt.Errorf("error reading cache file: %w", err)
	}

	for _, d := range dataleaks {
		fmt.Println(settings.Header.Render(d.Name))
		fmt.Println(settings.Base.Render("Length: ", fmt.Sprintf("%d", d.Length), " rows"))
		fmt.Println()
	}

	return nil
}
