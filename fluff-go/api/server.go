package api

import (
	"net/http"
	"time"

	"github.com/NoSoundLeR/fluff/fluff-go/api/handler"
	"github.com/NoSoundLeR/fluff/fluff-go/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Server ...
type Server struct {
	addr   string
	logger *zap.SugaredLogger
	router chi.Router
	db     *db.Database
}

// NewServer ...
func NewServer(bindAddr string, dbAddr string) *Server {
	config := zap.NewProductionConfig()
	config.Development = true
	config.Level.SetLevel(zapcore.InfoLevel)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	return &Server{
		addr:   bindAddr,
		logger: sugar,
		router: chi.NewRouter(),
		db:     db.NewDatabase(dbAddr, sugar),
	}
}

// Run ...
func (s *Server) Run() {
	defer s.db.Close()

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(time.Second * 30))

	s.router.Get("/{key}", handler.GetLink(s.db))
	s.router.Post("/api/links", handler.CreateLink(s.db))

	http.ListenAndServe(s.addr, s.router)
}
