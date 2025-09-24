package osint

import (
	"strings"
	"time"

	"github.com/anotherhadi/eleakxir/backend/server"

	recon_email "github.com/anotherhadi/github-recon/github-recon/email"
	recon_username "github.com/anotherhadi/github-recon/github-recon/username"
	github_recon_settings "github.com/anotherhadi/github-recon/settings"
)

type GithubResult struct {
	Duration time.Duration
	Error    string

	UsernameResult *recon_username.UsernameResult
	EmailResult    *recon_email.EmailResult
}

func Search(s *server.Server, queryText, column string) *GithubResult {
	if !s.Settings.GithubRecon {
		return nil
	}
	gr := GithubResult{}
	now := time.Now()
	settings := github_recon_settings.GetDefaultSettings()
	settings.Token = s.Settings.GithubToken
	settings.DeepScan = s.Settings.GithubDeepMode
	if settings.Token != "null" && strings.TrimSpace(settings.Token) != "" {
		settings.Client = settings.Client.WithAuthToken(settings.Token)
	}
	settings.Silent = true

	queryText = strings.TrimSpace(queryText)

	if column == "email" || strings.HasSuffix(column, "_email") ||
		column == "username" || strings.HasSuffix(column, "_username") ||
		column == "" || column == "all" {
		if isValidEmail(queryText) {
			settings.Target = queryText
			settings.TargetType = github_recon_settings.TargetEmail
			result := recon_email.Email(settings)
			gr.EmailResult = &result
		} else if isValidUsername(queryText) {
			settings.Target = queryText
			settings.TargetType = github_recon_settings.TargetUsername
			result, err := recon_username.Username(settings)
			if err != nil {
				gr.Error = err.Error()
			}
			if result.User.Username == "" {
				gr.UsernameResult = nil
			} else {
				gr.UsernameResult = &result
			}
		} else {
			return nil
		}
	} else {
		return nil
	}

	gr.Duration = time.Since(now)
	return &gr
}

func isValidEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	if strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return false
	}
	if strings.Contains(email, " ") {
		return false
	}
	return true
}

func isValidUsername(username string) bool {
	if len(username) < 1 || len(username) > 39 {
		return false
	}
	if strings.Contains(username, " ") {
		return false
	}
	return true
}
