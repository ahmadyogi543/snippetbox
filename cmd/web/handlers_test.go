package main

import (
	"net/http"
	"testing"

	"github.com/ahmadyogi543/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	code, _, body := server.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
