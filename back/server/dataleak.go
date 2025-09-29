package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

type Dataleak struct {
	Path    string
	Name    string
	Columns []string
	Length  uint64
	Size    uint64
	ModTime time.Time
}

const CACHE_FILENAME = "dataleaks_cache.json"

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

	dataleakMap := make(map[string]Dataleak, len(dataleaks))
	for _, d := range dataleaks {
		dataleakMap[d.Path] = d
	}

	filteredDataleaks := []Dataleak{}
	writeOutput := false
	parquetFiles := getAllParquetFiles(s.Settings.Folders)

	for _, p := range parquetFiles {
		currentModTime := getModTime(p)
		cachedDataleak, found := dataleakMap[p]

		if found {
			if currentModTime.Equal(cachedDataleak.ModTime) {
				filteredDataleaks = append(filteredDataleaks, cachedDataleak)
			} else {
				log.Info("File modification time changed, re-caching dataleak", "path", p)
				writeOutput = true
				filteredDataleaks = append(filteredDataleaks, getDataleak(*s, p))
			}
			delete(dataleakMap, p)
		} else {
			log.Info("Found new dataleak file, caching", "path", p)
			writeOutput = true
			filteredDataleaks = append(filteredDataleaks, getDataleak(*s, p))
		}
	}

	for path := range dataleakMap {
		log.Info("Removing non-existent file from cache", "path", path)
		writeOutput = true
	}

	dataleaks = filteredDataleaks

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
		ModTime: getModTime(path),
	}
}
