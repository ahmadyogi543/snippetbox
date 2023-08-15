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

func TestSnippetView(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	tests := []struct {
		name         string
		urlPath      string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid ID",
			urlPath:      "/snippet/view/1",
			expectedCode: http.StatusOK,
			expectedBody: "This is a content inside the mock snippet.",
		},
		{
			name:         "Non-existent ID",
			urlPath:      "/snippet/view/1000",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Negative ID",
			urlPath:      "/snippet/view/-1",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Decimal ID",
			urlPath:      "/snippet/view/1.45",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "String ID",
			urlPath:      "/snippet/view/abc",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Empty ID",
			urlPath:      "/snippet/view/",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := server.get(t, test.urlPath)

			assert.Equal(t, code, test.expectedCode)
			if test.expectedBody != "" {
				assert.StringContaints(t, body, test.expectedBody)
			}
		})
	}
}
