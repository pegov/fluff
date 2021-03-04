package handler

import (
	"net/http"

	"github.com/NoSoundLeR/fluff/fluff-go/db"
	"github.com/go-chi/chi"
)

// GetLink ...
func GetLink(db db.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key == "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		url, err := db.GetLink(key)
		if err != nil {
			http.Error(w, "404 page not found", 404)
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
}
