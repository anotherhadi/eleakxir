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
	Inactive bool

	UsernameResult *recon_username.UsernameResult
	EmailResult    *recon_email.EmailResult
}

func GithubSearch(s *server.Server, queryText, queryType string) GithubResult {
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

	if queryType == "email" {
		settings.Target = queryText
		settings.TargetType = github_recon_settings.TargetEmail
		result := recon_email.Email(settings)
		gr.EmailResult = &result
	} else if queryType == "username" {
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
		return GithubResult{Inactive: true}
	}

	gr.Duration = time.Since(now)
	return gr
}
