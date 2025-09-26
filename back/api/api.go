package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/anotherhadi/eleakxir/backend/search"
	"github.com/anotherhadi/eleakxir/backend/server"
	"github.com/gin-gonic/gin"
)

func searchWorker(s *server.Server, cache *map[string]*search.Result, searchQueue chan string) {
	for id := range searchQueue {
		s.Mu.RLock()
		r, exists := (*cache)[id]
		s.Mu.RUnlock()

		if !exists {
			continue
		}

		search.Search(s, r.Query, r, s.Mu)
	}
}

func routes(s *server.Server, cache *map[string]*search.Result, searchQueue chan string) {
	s.Router.Use(
		func(c *gin.Context) {
			if s.Settings.Password != "" {
				password := c.GetHeader("X-Password")
				if password != s.Settings.Password {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
					return
				}
			}
			c.Next()
		},
	)

	s.Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Settings":       s.Settings,
			"Dataleaks":      s.Dataleaks,
			"TotalDataleaks": s.TotalDataleaks,
			"TotalRows":      s.TotalRows,
			"TotalSize":      s.TotalSize,
		})
	})

	s.Router.GET("/history", func(c *gin.Context) {
		type historyItem struct {
			Id      string
			Status  string
			Date    time.Time
			Query   search.Query
			Results int
		}
		var history []historyItem
		s.Mu.RLock()
		for _, r := range *cache {
			resultsCount := 0
			if r.LeakResult.Rows != nil {
				resultsCount = len(r.LeakResult.Rows)
			}
			history = append(history, historyItem{
				Id:      r.Id,
				Status:  r.Status,
				Date:    r.Date,
				Query:   r.Query,
				Results: resultsCount,
			})
		}
		s.Mu.RUnlock()
		for i := 0; i < len(history)-1; i++ {
			for j := 0; j < len(history)-i-1; j++ {
				if history[j].Date.Before(history[j+1].Date) {
					history[j], history[j+1] = history[j+1], history[j]
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"History": history,
		})
	})

	s.Router.POST("/search", func(c *gin.Context) {
		var query search.Query
		if err := c.BindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "invalid JSON"})
			return
		}
		query = cleanQuery(query)
		if len(query.Text) <= s.Settings.MinimumQueryLength {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "query too short"})
			return
		}
		id := search.EncodeQueryID(query, *s.TotalDataleaks)
		s.Mu.RLock()
		_, exists := (*cache)[id]
		s.Mu.RUnlock()
		if exists {
			c.JSON(http.StatusOK, gin.H{"Id": id})
			return
		}
		r := search.Result{
			Id:     id,
			Date:   time.Now(),
			Status: "queued",
			Query:  query,
		}

		s.Mu.Lock()
		(*cache)[id] = &r
		s.Mu.Unlock()

		searchQueue <- id

		c.JSON(http.StatusOK, gin.H{"Id": id})
	})

	s.Router.GET("/search/:id", func(c *gin.Context) {
		id := c.Param("id")
		s.Mu.RLock()
		r, exists := (*cache)[id]
		s.Mu.RUnlock()
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"Error": "not found"})
			return
		}
		c.JSON(http.StatusOK, r)
	})
}

func Init(s *server.Server) {
	if !s.Settings.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	s.Router = gin.Default()
	s.Router.Use(CORSMiddleware())

	cache := make(map[string]*search.Result)
	searchQueue := make(chan string, 100)

	go func() {
		for {
			time.Sleep(time.Minute)
			deleteOldCache(s, &cache)
		}
	}()

	go searchWorker(s, &cache, searchQueue)

	routes(s, &cache, searchQueue)
}

func deleteOldCache(s *server.Server, cache *map[string]*search.Result) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	now := time.Now()
	for id, r := range *cache {
		if now.Sub(r.Date) > s.Settings.MaxCacheDuration {
			delete(*cache, id)
		}
	}
}

func cleanQuery(q search.Query) search.Query {
	q.Column = strings.ToLower(strings.TrimSpace(q.Column))
	q.Column = strings.Join(strings.Fields(q.Column), " ")
	q.Column = strings.ReplaceAll(q.Column, "`", "")
	q.Column = strings.ReplaceAll(q.Column, "'", "")
	q.Column = strings.ReplaceAll(q.Column, "-", "_")
	q.Column = strings.ReplaceAll(q.Column, " ", "_")
	q.Column = strings.ReplaceAll(q.Column, "\"", "")

	q.Text = strings.TrimSpace(q.Text)
	q.Text = strings.Join(strings.Fields(q.Text), " ")

	return q
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "X-Password, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
