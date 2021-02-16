package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NoSoundLeR/fluff/fluff-go/db"
	"github.com/NoSoundLeR/fluff/fluff-go/link"
)

// CreateLink ...
func CreateLink(db db.Setter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var link link.Link
		if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(link.Key) == 0 {
			link.Key = db.GetKey()
		}
		if ok, err := link.ValidateURL(); !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		db.SetLink(link)
		if err := json.NewEncoder(w).Encode(link); err != nil {
			http.Error(w, "can't encode", http.StatusBadRequest)
		}
	}
}
