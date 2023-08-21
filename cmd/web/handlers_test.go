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
				assert.StringContains(t, body, test.expectedBody)
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
		form.Add("email", "ayogi@snippetbox.sh")
		form.Add("password", "12345678")
		form.Add("csrf_token", csrfToken)
		server.postForm(t, "/user/login", form)

		code, _, body := server.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<form action="/snippet/create" method="POST">`)
	})
}

func TestSnippetCreatePost(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	tests := []struct {
		name         string
		title        string
		content      string
		expires      string
		expectedCode int
	}{
		{
			name:         "Valid Form",
			title:        "A Title",
			content:      "This is a content example",
			expires:      "365",
			expectedCode: http.StatusSeeOther,
		},
		{
			name:         "Empty Field",
			title:        "",
			content:      "",
			expires:      "7",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "Invalid Expires",
			title:        "A Title",
			content:      "This is a content example",
			expires:      "1000",
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, body := server.get(t, "/user/login")
			csrfToken := extractCSRFToken(t, body)
			form := url.Values{}
			form.Add("name", "Ahmad Yogi")
			form.Add("email", "ayogi@snippetbox.sh")
			form.Add("password", "12345678")
			form.Add("csrf_token", csrfToken)
			server.postForm(t, "/user/login", form)

			_, _, body = server.get(t, "/snippet/create")
			csrfToken = extractCSRFToken(t, body)
			form = url.Values{}
			form.Add("title", test.title)
			form.Add("content", test.content)
			form.Add("expires", test.expires)
			form.Add("csrf_token", csrfToken)

			code, _, _ := server.postForm(t, "/snippet/create", form)
			assert.Equal(t, code, test.expectedCode)
		})
	}
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
				assert.StringContains(t, body, test.expectedFormTag)
			}
		})
	}
}

func TestAccountView(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	_, _, body := server.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("email", "ayogi@snippetbox.sh")
	form.Add("password", "12345678")
	form.Add("csrf_token", csrfToken)

	server.postForm(t, "/user/login", form)

	code, _, _ := server.get(t, "/account/view")
	assert.Equal(t, code, http.StatusOK)
}

func TestPasswordUpdate(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := server.get(t, "/account/password/update")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := server.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "ayogi@snippetbox.sh")
		form.Add("password", "12345678")
		form.Add("csrf_token", csrfToken)
		server.postForm(t, "/user/login", form)

		code, _, _ := server.get(t, "/account/password/update")
		assert.Equal(t, code, http.StatusOK)
	})
}

func TestPasswordUpdatePost(t *testing.T) {
	app := newTestApp(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	tests := []struct {
		name               string
		currentPassword    string
		newPassword        string
		confirmNewPassword string
		expectedCode       int
	}{
		{
			name:               "Valid Form",
			currentPassword:    "12345678",
			newPassword:        "87654321",
			confirmNewPassword: "87654321",
			expectedCode:       http.StatusSeeOther,
		},
		{
			name:               "Empty Form",
			currentPassword:    "",
			newPassword:        "",
			confirmNewPassword: "",
			expectedCode:       http.StatusUnprocessableEntity,
		},
		{
			name:               "Not Match Confirm New Password",
			currentPassword:    "12345678",
			newPassword:        "87654321",
			confirmNewPassword: "12345678",
			expectedCode:       http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, body := server.get(t, "/user/login")
			csrfToken := extractCSRFToken(t, body)
			form := url.Values{}
			form.Add("name", "Ahmad Yogi")
			form.Add("email", "ayogi@snippetbox.sh")
			form.Add("password", "12345678")
			form.Add("csrf_token", csrfToken)
			server.postForm(t, "/user/login", form)

			_, _, body = server.get(t, "/account/password/update")
			csrfToken = extractCSRFToken(t, body)
			form = url.Values{}
			form.Add("current_password", test.currentPassword)
			form.Add("new_password", test.newPassword)
			form.Add("confirm_new_password", test.confirmNewPassword)
			form.Add("csrf_token", csrfToken)

			code, _, _ := server.postForm(t, "/account/password/update", form)
			assert.Equal(t, code, test.expectedCode)
		})
	}
}
