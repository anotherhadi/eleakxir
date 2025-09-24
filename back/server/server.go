package server

import (
	"database/sql"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	_ "github.com/marcboeker/go-duckdb"
)

type Server struct {
	Settings ServerSettings

	Dataleaks *[]Dataleak

	TotalRows      *uint64
	TotalDataleaks *uint64
	TotalSize      *uint64 // MB

	Router *gin.Engine
	Duckdb *sql.DB
	Mu     *sync.RWMutex
}

func NewServer() *Server {
	zero := uint64(0)
	emptyDataleak := []Dataleak{}
	s := &Server{
		Settings:       LoadServerSettings(),
		Mu:             &sync.RWMutex{},
		TotalDataleaks: &zero,
		TotalRows:      &zero,
		TotalSize:      &zero,
		Dataleaks:      &emptyDataleak,
	}

	var err error

	s.Duckdb, err = sql.Open("duckdb", "")
	if err != nil {
		panic(err)
	}

	err = Cache(s)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(s.Settings.ReloadDataleaksInterval)
			err := Cache(s)
			if err != nil {
				log.Error(err)
			}
		}
	}()

	return s
}
