package main

import (
	"net/http"
	"net/url"
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

func TestSnippetCreate(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := server.get(t, "/snippet/create")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := server.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("name", "Ahmad Yogi")
		form.Add("email", "ayogi@snippetbox.sh")
		form.Add("password", "12345678")
		form.Add("csrf_token", csrfToken)
		server.postForm(t, "/user/login", form)

		code, _, body := server.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContaints(t, body, `<form action="/snippet/create" method="POST">`)
	})
}

func TestUserSignup(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	_, _, body := server.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "Ahmad Yogi"
		validPassword = "12345678"
		validEmail    = "ayogi@snippetbox.sh"
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		name            string
		userName        string
		userEmail       string
		userPassword    string
		csrfToken       string
		expectedCode    int
		expectedFormTag string
	}{
		{
			name:         "Valid Submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			expectedCode: http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "invalid csrf token",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:            "Empty Name",
			userName:        "",
			userEmail:       validEmail,
			userPassword:    validPassword,
			csrfToken:       validCSRFToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Empty Email",
			userName:        validName,
			userEmail:       "",
			userPassword:    validPassword,
			csrfToken:       validCSRFToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Empty Password",
			userName:        validName,
			userEmail:       validEmail,
			userPassword:    "",
			csrfToken:       validCSRFToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Invalid Email",
			userName:        validName,
			userEmail:       "ambatu@kam.",
			userPassword:    validPassword,
			csrfToken:       validCSRFToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Short Password",
			userName:        validName,
			userEmail:       validEmail,
			userPassword:    "123",
			csrfToken:       validCSRFToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Duplicate Email",
			userName:        validName,
			userEmail:       "duplicate@snippetbox.sh",
			userPassword:    validPassword,
			csrfToken:       validCSRFToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)

			code, _, body := server.postForm(t, "/user/signup", form)
			assert.Equal(t, code, test.expectedCode)
			if test.expectedFormTag != "" {
				assert.StringContaints(t, body, test.expectedFormTag)
			}
		})
	}
}
