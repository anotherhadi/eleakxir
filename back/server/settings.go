package server

import (
	"os"
	"strconv"
	"strings"
	"time"

	github_recon_settings "github.com/anotherhadi/github-recon/settings"
)

type ServerSettings struct {
	Port               int `json:"-"` // Port to run the server on
	Debug              bool
	Password           string `json:"-"` // Do not expose the password in JSON
	MinimumQueryLength int
	MaxCacheDuration   time.Duration // Delete a search from the cache after this duration

	// Dataleaks
	Folders                 []string // Folders to search in for parquets, recursive
	CacheFolder             string
	BaseColumns             []string      // Use these columns when column="all"
	Limit                   int           // Limit number of rows returned
	ReloadDataleaksInterval time.Duration // Reload dataleaks files from disk every X

	// OSINT Tools
	GithubRecon       bool   // Activate github-recon OSINT tool
	GithubToken       string `json:"-"` // Github token for github-recon
	GithubTokenLoaded bool
	GithubDeepMode    bool // Deep mode for github-recon
}

func LoadServerSettings() ServerSettings {
	ss := ServerSettings{
		Port:               getEnvPortOrDefault("PORT", 9198),
		Debug:              getEnvBoolOrDefault("DEBUG", false),
		Password:           getEnvStringOrDefault("PASSWORD", ""),
		MinimumQueryLength: getEnvIntOrDefault("MINIMUM_QUERY_LENGTH", 3),
		MaxCacheDuration:   getEnvDurationOrDefault("MAX_CACHE_DURATION", 24*time.Hour),

		// Dataleaks
		Folders:                 getEnvStringListOrDefault("DATALEAKS_FOLDERS", []string{}),
		CacheFolder:             getEnvStringOrDefault("DATALEAKS_CACHE_FOLDER", ""),
		BaseColumns:             getEnvStringListOrDefault("BASE_COLUMNS", []string{"email", "username", "password", "full_name", "phone", "url"}),
		Limit:                   getEnvIntOrDefault("LIMIT", 200),
		ReloadDataleaksInterval: getEnvDurationOrDefault("RELOAD_DATALEAKS_INTERVAL", 20*time.Minute),

		// OSINT Tools
		GithubRecon:    getEnvBoolOrDefault("GITHUB_RECON", true),
		GithubToken:    getEnvStringOrDefault("GITHUB_TOKEN", "null"),
		GithubDeepMode: getEnvBoolOrDefault("GITHUB_DEEP_MODE", false),
	}

	if ss.GithubToken == "null" || strings.TrimSpace(ss.GithubToken) == "" {
		ss.GithubToken = github_recon_settings.GetToken()
	}

	if ss.GithubToken != "null" && strings.TrimSpace(ss.GithubToken) != "" {
		ss.GithubTokenLoaded = true
	}

	return ss
}

func getEnvStringOrDefault(envKey, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(envKey))
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBoolOrDefault(envKey string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(envKey))
	if value == "" {
		return defaultValue
	}
	value = strings.ToLower(value)
	if value == "true" || value == "1" {
		return true
	} else if value == "false" || value == "0" {
		return false
	}
	return defaultValue
}

func getEnvDurationOrDefault(envKey string, defaultValue time.Duration) time.Duration {
	v := getEnvStringOrDefault(envKey, "")
	if v == "" {
		return defaultValue
	}
	t, err := time.ParseDuration(v)
	if err != nil {
		return defaultValue
	}
	return t
}

func getEnvStringListOrDefault(envKey string, defaultValue []string) []string {
	value := strings.TrimSpace(os.Getenv(envKey))
	if value == "" {
		return defaultValue
	}
	l := strings.Split(value, ",")
	for i := range l {
		l[i] = strings.TrimSpace(l[i])
	}
	return l
}

func getEnvIntOrDefault(envKey string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(envKey))
	if value == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func getEnvPortOrDefault(envKey string, defaultValue int) int {
	p := getEnvIntOrDefault(envKey, defaultValue)
	if p <= 0 || p >= 65534 {
		return defaultValue
	}
	return p
}
