package main

import (
	"net/http"
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
