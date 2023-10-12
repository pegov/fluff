package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksHandlers(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/links", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
