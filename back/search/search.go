package search

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
}

type Result struct {
	Id     string
	Date   time.Time
	Status string // "pending", "completed"
	Query  Query

	LeakResult   dataleak.LeakResult
	GithubResult osint.GithubResult
}

func Search(s *server.Server, q Query, r *Result, mu *sync.RWMutex) {
	var wg sync.WaitGroup

	mu.Lock()
	r.Date = time.Now()
	r.Status = "pending"
	r.Query = q
	mu.Unlock()

	wg.Add(2)

	go func() {
		leakResult := dataleak.Search(s, q.Text, q.Column, q.ExactMatch)
		mu.Lock()
		r.LeakResult = leakResult
		mu.Unlock()
		wg.Done()
	}()

	go func() {
		githubResult := osint.Search(s, q.Text, q.Column)
		mu.Lock()
		r.GithubResult = *githubResult
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
