package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ahmadyogi543/snippetbox/internal/mocks"
	"github.com/alexedwards/scs/v2"
)

type testServer struct {
	*httptest.Server
}

func newTestApp(t *testing.T) *App {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &App{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		templateCache:  templateCache,
		sessionManager: sessionManager,
		users:          &mocks.UserModel{},
		snippets:       &mocks.SnippetModel{},
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	server := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	server.Client().Jar = jar

	server.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{server}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	result, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return result.StatusCode, result.Header, string(body)
}
