package osint

import (
	"time"

	"github.com/anotherhadi/eleakxir/backend/server"
	gravatar_recon "github.com/anotherhadi/gravatar-recon"
)

type GravatarResult struct {
	Duration time.Duration
	Error    string
	Inactive bool

	Results []gravatar_recon.GravatarProfile
}

func GravatarSearch(s *server.Server, queryText string) GravatarResult {
	gr := GravatarResult{}
	now := time.Now()
	results, err := gravatar_recon.GetGravatarProfiles(queryText)

	if err != nil {
		gr.Error = err.Error()
		return gr
	}

	gr.Results = *results

	gr.Duration = time.Since(now)
	return gr
}
