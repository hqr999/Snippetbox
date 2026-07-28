package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hqr999/Snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {

	// Initialize the new httptest.ResponseRecorder and a dummt http.Request
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our commonHeaders
	// Middleware, which writes a 200 status code and an "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to our commonHeaders middleware. Because
	// commonHeaders *returns* an http.Handler we can still its ServeHTTP()
	// method, passing in the http.ResponseRecorder and dummy http.Request to
	// execute it.
	commonHeaders(next).ServeHTTP(rr, req)

	// Call the Result() method on the http.ResponseRecorder to get the results
	// of the test.
	res := rr.Result()
	defer res.Body.Close()

	// Check that the middleware has correctly set the Content-Secutiry-Policy
	// header on the response.
	expecVal := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), expecVal)

	// Check that the middleware has correctly set the Referrer-Policy 
	// header on the response 
	expecVal = "origin-when-cross-origin"
	assert.Equal(t,res.Header.Get("Referrer-Policy"),expecVal)


	// Check that the middleware has correctly set the X-Content-Type-Options
	// header on the response.
	expecVal = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), expecVal)

	// Check that the middleware has correctly set the X-Frame-Options header
	// on the response.
	expecVal = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), expecVal)

	// Check that the middleware has correctly set the X-XSS-Protection header 
	// on the response. 
	expecVal = "0"
	assert.Equal(t,res.Header.Get("X-XSS-Protection"),expecVal)


	// Check that the middleware has correctly set the server Header on the
	// response.
	expecVal = "Go"
	assert.Equal(t, res.Header.Get("Server"), expecVal)

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}
