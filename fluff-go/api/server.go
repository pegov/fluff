package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NoSoundLeR/fluff/fluff-go/db"
	"github.com/NoSoundLeR/fluff/fluff-go/link"
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

func (s *Server) getLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key == "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		url, err := s.db.GetLink(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
}

func (s *Server) createLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var link link.Link
		if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var custom bool
		if len(link.Key) == 0 {
			link.Key = s.db.GetFreeKey()
		} else {
			custom = true
		}
		if ok, err := link.ValidateURL(); !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.db.SetLink(link, custom)
		if err := json.NewEncoder(w).Encode(link); err != nil {
			http.Error(w, "can't encode", http.StatusBadRequest)
		}
	}
}

// Run ...
func (s *Server) Run() {
	s.logger.Info("Adding middlewares...")
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(time.Second * 30))

	s.logger.Info("Adding router...")
	s.router.Get("/{key}", s.getLink())
	s.router.Post("/api/links", s.createLink())
	s.logger.Infof("Listening on addr: %s", s.addr)
	http.ListenAndServe(s.addr, s.router)
	defer s.db.Close()
}
