package misc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
)

func flattenJSON(prefix string, in map[string]any, out map[string]any) {
	for k, v := range in {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch child := v.(type) {
		case map[string]any:
			flattenJSON(key, child, out)
		case []any:
			if len(child) == 0 {
				out[key] = ""
			} else {
				for i, item := range child {
					tempMap := make(map[string]any)
					subKey := fmt.Sprintf("%d", i)

					if obj, ok := item.(map[string]any); ok {
						flattenJSON(key+"."+subKey, obj, out)
					} else if arr, ok := item.([]any); ok {
						tempMap[subKey] = arr
						flattenJSON(key, tempMap, out)
					} else {
						out[key+"."+subKey] = fmt.Sprintf("%v", item)
					}
				}
			}
		default:
			out[key] = fmt.Sprintf("%v", child)
		}
	}
}

func flattenJSONFile(inputFile string, outputFile string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") {
			var arr []map[string]any
			if err := json.Unmarshal([]byte(line), &arr); err != nil {
				return fmt.Errorf("invalid JSON array: %w", err)
			}
			for _, obj := range arr {
				flat := make(map[string]any)
				flattenJSON("", obj, flat)
				b, err := json.Marshal(flat)
				if err != nil {
					return fmt.Errorf("failed to marshal flattened JSON: %w", err)
				}
				writer.Write(b)
				writer.WriteString("\n")
			}
		} else {
			var obj map[string]any
			if err := json.Unmarshal([]byte(line), &obj); err != nil {
				return fmt.Errorf("invalid JSON object: %w", err)
			}
			flat := make(map[string]any)
			flattenJSON("", obj, flat)
			b, err := json.Marshal(flat)
			if err != nil {
				return fmt.Errorf("failed to marshal flattened JSON: %w", err)
			}
			writer.Write(b)
			writer.WriteString("\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func JsonToParquet(lu settings.LeakUtils, inputFile string, outputFile string) error {
	tmpFile := filepath.Join("/tmp", "leak-utils.flat.json")
	err := flattenJSONFile(inputFile, tmpFile)
	defer os.Remove(tmpFile)

	query := fmt.Sprintf(`COPY (FROM read_json('%s', union_by_name=true)) TO '%s' (FORMAT 'parquet', COMPRESSION '%s', ROW_GROUP_SIZE 200000);`, tmpFile, outputFile, lu.Compression)

	_, err = lu.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("duckdb copy error: %w", err)
	}

	return nil
}
