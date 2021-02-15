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

type Server struct {
	addr   string
	logger *zap.SugaredLogger
	router chi.Router
	db     *db.Database
}

// NewServer ...
func NewServer(addr string) *Server {
	config := zap.NewProductionConfig()
	config.Development = true
	config.Level.SetLevel(zapcore.InfoLevel)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	return &Server{
		addr:   addr,
		logger: sugar,
		router: chi.NewRouter(),
		db:     db.NewDatabase(sugar),
	}
}

// Run ...
func (s *Server) Run() {
	s.logger.Info("Adding middlewares...")
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(time.Second * 30))

	s.logger.Info("Adding router...")
	s.router.Get("/{key}", handler.GetLink(s.db))
	s.router.Post("/api/links", handler.CreateLink(s.db))
	s.logger.Infof("Listening on addr: %s", s.addr)
	http.ListenAndServe(s.addr, s.router)
	defer s.db.Close()
}
