package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type Dataleak struct {
	Path    string
	Name    string
	Columns []string
	Length  uint64
	Size    uint64
}

const CACHE_FILENAME = "dataleaks_cache.json"

// TODO: check os.FileInfo.ModTime() to see if the file has changed since last cache update

func Cache(s *Server) error {
	if len(s.Settings.Folders) == 0 {
		return nil
	}
	if s.Settings.CacheFolder == "" {
		s.Settings.CacheFolder = s.Settings.Folders[0]
	}
	if err := createDirectoryIfNotExists(s.Settings.CacheFolder); err != nil {
		return err
	}

	cacheFile := filepath.Join(s.Settings.CacheFolder, CACHE_FILENAME)
	dataleaks := []Dataleak{}

	data, err := os.ReadFile(cacheFile)
	if err == nil {
		if err := json.Unmarshal(data, &dataleaks); err != nil {
			log.Warn("Failed to unmarshal dataleaks cache", "error", err)
		}
	} else {
		log.Warn("Failed to read dataleaks cache file", "error", err)
	}

	// Filter out non-existent files
	filteredDataleaks := []Dataleak{}
	writeOutput := false
	for _, d := range dataleaks {
		if _, err := os.Stat(d.Path); err == nil {
			filteredDataleaks = append(filteredDataleaks, d)
		} else if os.IsNotExist(err) {
			log.Info("Removing non-existent file from cache", "path", d.Path)
			writeOutput = true
		} else {
			log.Error("Error checking file existence", "path", d.Path, "error", err)
		}
	}
	dataleaks = filteredDataleaks

	// Create a map for quick lookups
	dataleakMap := make(map[string]struct{}, len(dataleaks))
	for _, d := range dataleaks {
		dataleakMap[d.Path] = struct{}{}
	}

	// Add new files
	parquetFiles := getAllParquetFiles(s.Settings.Folders)
	for _, p := range parquetFiles {
		if _, found := dataleakMap[p]; found {
			continue
		}

		writeOutput = true
		dataleaks = append(dataleaks, getDataleak(*s, p))
	}

	if writeOutput {
		data, err := json.MarshalIndent(dataleaks, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling cache: %w", err)
		}
		if err := os.WriteFile(cacheFile, data, 0644); err != nil {
			return fmt.Errorf("error writing cache: %w", err)
		}
	}

	s.Dataleaks = &dataleaks
	totalDataleaks := uint64(len(dataleaks))
	totalRows := uint64(0)
	totalSize := uint64(0)
	for _, d := range dataleaks {
		totalRows += d.Length
		totalSize += d.Size
	}
	s.TotalDataleaks = &totalDataleaks
	s.TotalSize = &totalSize
	s.TotalRows = &totalRows
	return nil
}

func getDataleak(s Server, path string) Dataleak {
	return Dataleak{
		Path:    path,
		Name:    FormatParquetName(path),
		Columns: getParquetColumns(s, path),
		Length:  getParquetLength(s, path),
		Size:    getFileSize(path),
	}
}
