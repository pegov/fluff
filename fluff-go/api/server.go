package api

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/NoSoundLeR/fluff/fluff-go/api/handler"
	"github.com/NoSoundLeR/fluff/fluff-go/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server ...
type Server struct {
	addr   string
	router chi.Router
	db     *db.Database
}

// NewServer ...
func NewServer() *Server {
	bindAddr, ok := os.LookupEnv("BASE_URL")
	if !ok {
		bindAddr = "0.0.0.0:8000"
	}
	redisAddr, ok := os.LookupEnv("REDIS_URL")
	if ok {
		// redis://127.0.0.1:6379/ -> 127.0.0.1:6379
		redisAddr = strings.ReplaceAll(redisAddr, "redis://", "")
		redisAddr = strings.ReplaceAll(redisAddr, "/", "")
	} else {
		redisAddr = "127.0.0.1:6379"
	}
	return &Server{
		addr:   bindAddr,
		router: chi.NewRouter(),
		db:     db.NewDatabase(redisAddr),
	}
}

// Run ...
func (s *Server) Run() {
	defer s.db.Close()

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/{key}", handler.GetLink(s.db))
	s.router.Post("/api/links", handler.CreateLink(s.db))

	err := http.ListenAndServe(s.addr, s.router)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
