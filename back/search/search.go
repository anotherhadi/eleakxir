package search

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/anotherhadi/eleakxir/backend/search/dataleak"
	"github.com/anotherhadi/eleakxir/backend/search/osint"
	"github.com/anotherhadi/eleakxir/backend/server"
)

type Query struct {
	Text       string
	Column     string // The column to search in (e.g., "email", "password", etc.
	ExactMatch bool   // Whether to search for an exact match

	// Services
	Datawells     bool // Whether to include datawells in the search
	GithubRecon   bool // Whether to include github-recon in the search
	GravatarRecon bool // Whether to include gravatar-recon in the search
}

type Result struct {
	Id     string
	Date   time.Time
	Status string // "pending", "completed"
	Query  Query

	LeakResult     dataleak.LeakResult
	GithubResult   osint.GithubResult
	GravatarResult osint.GravatarResult
}

func Search(s *server.Server, q Query, r *Result, mu *sync.RWMutex) {
	var wg sync.WaitGroup

	mu.Lock()
	r.Date = time.Now()
	r.Status = "pending"
	r.Query = q
	mu.Unlock()

	wg.Add(3)
	go func() {
		if !q.Datawells {
			mu.Lock()
			r.LeakResult = dataleak.LeakResult{Inactive: true}
			mu.Unlock()
			wg.Done()
			return
		}
		leakResult := dataleak.Search(s, q.Text, q.Column, q.ExactMatch)
		mu.Lock()
		r.LeakResult = leakResult
		mu.Unlock()
		wg.Done()
	}()

	cleanQueryText := strings.TrimPrefix(q.Text, "^")
	cleanQueryText = strings.TrimSuffix(q.Text, "$")
	isEmail := false
	isUsername := false

	if q.Column == "email" || strings.HasSuffix(q.Column, "_email") ||
		q.Column == "username" || strings.HasSuffix(q.Column, "_username") ||
		q.Column == "" || q.Column == "all" {
		if isValidEmail(cleanQueryText) {
			isEmail = true
		} else if isValidUsername(cleanQueryText) {
			isUsername = true
		}
	}

	go func() {
		if !q.GithubRecon || !s.Settings.GithubRecon || (!isEmail && !isUsername) {
			mu.Lock()
			r.GithubResult = osint.GithubResult{Inactive: true}
			mu.Unlock()
			wg.Done()
			return
		}
		var githubResult osint.GithubResult
		if isEmail {
			githubResult = osint.GithubSearch(s, cleanQueryText, "email")
		} else if isUsername {
			githubResult = osint.GithubSearch(s, cleanQueryText, "username")
		}
		mu.Lock()
		r.GithubResult = githubResult
		mu.Unlock()
		wg.Done()
	}()

	go func() {
		if !q.GravatarRecon || !s.Settings.GravatarRecon || !isEmail {
			mu.Lock()
			r.GravatarResult = osint.GravatarResult{Inactive: true}
			mu.Unlock()
			wg.Done()
			return
		}
		gravatarResult := osint.GravatarSearch(s, cleanQueryText)
		mu.Lock()
		r.GravatarResult = gravatarResult
		mu.Unlock()
		wg.Done()
	}()

	wg.Wait()

	mu.Lock()
	r.Status = "completed"
	mu.Unlock()
}

func EncodeQueryID(q Query, dataleaksCount uint64) string {
	raw, _ := json.Marshal(q)
	return fmt.Sprintf("%d:%s", dataleaksCount, base64.URLEncoding.EncodeToString(raw))
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
