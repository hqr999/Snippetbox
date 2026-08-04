package main

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hqr999/Snippetbox/internal/assert"
)

func TestSnippetView(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table driven tests in order to check the responses by our
	// application so far.
	tests := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid ID",
			urlPath:    "/snippet/view/1",
			wantStatus: http.StatusOK,
			wantBody:   "An old silent pond..",
		},
		{
			name:       "Non-existent ID",
			urlPath:    "/snippet/view/2",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Negative ID",
			urlPath:    "/snippet/view/-1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Decimal ID",
			urlPath:    "/snippet/view/1.56",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "String ID",
			urlPath:    "/snippet/view/ab",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Empty ID",
			urlPath:    "/snippet/view/",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			// Reset the test server client cookie jar at the start of each
			// sub-test, so that any cookies set in previous tests are removed
			// and don´t affect this test.
			ts.resetClientCookieJar(t)

			response := ts.get(t, v.urlPath)
			// Use the assert.Equal() to check the response status, and the
			// assert.True() function in conjunction with strings.Contains() to
			// make sure that the response body contains the expected content.
			assert.Equal(t, response.status, v.wantStatus)
			assert.True(t, strings.Contains(response.body, v.wantBody))

		})
	}

}

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	t_server := newTestServer(t, app.routes())
	defer t_server.Close()

	res := t_server.get(t, "/ping")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}

func TestUserSignup(t *testing.T) {
	// Create the application struct containing our mocked dependencies and
	// establish a new test server.
	app := newTestApplication(t)
	test_server := newTestServer(t, app.routes())
	defer test_server.Close()

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = `<form action="/user/signup" method="post" novalidate`
	)

	tests := []struct {
		name              string
		userName          string
		userEmail         string
		userPassword      string
		useValidCSRFToken bool
		wantStatus        int
		wantFormTag       string
	}{
		{
			name:              "Valid Submission",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusSeeOther,
		},
		{
			name:              "Invalid CSRF Token",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: false,
			wantStatus:        http.StatusBadRequest,
		},
		{
			name:              "Empty Name",
			userName:          "",
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty Email",
			userName:          validName,
			userEmail:         "",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty Password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Invalid Email",
			userName:          validName,
			userEmail:         "bob@example.",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Short Password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "pa$$",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Duplicate Email",
			userName:          validName,
			userEmail:         "dupe@example.com",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			// Reset the cookie jar for each sub-test.
			test_server.resetClientCookieJar(t)

			// Make a GET /user/signup request. This will automatically
			// add the CSRF cookie from the response to the test clients cookie
			// jar, and we can extract the CSRF token from the response body.
			res := test_server.get(t, "/user/signup")

			// Build up the form values for the sub-test, including the CSRF
			// token if appropriate.
			form := url.Values{}
			form.Add("name", v.userName)
			form.Add("email", v.userEmail)
			form.Add("password", v.userPassword)
			if v.useValidCSRFToken {
				form.Add("csrf_token", extractCSRFToken(t, res.body))
			}

			// Make the POST /user/signup request using the form values we
			// created above. The request will automatically include the CSRF
			// cookie from the test client´s cookie jar.
			res = test_server.postForm(t, "/user/signup", form)

			// And finally, test the response data.
			assert.Equal(t, res.status, v.wantStatus)
			assert.True(t, strings.Contains(res.body, v.wantFormTag))

		})
	}
}
