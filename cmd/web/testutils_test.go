package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{logger: slog.New(slog.DiscardHandler)}

}

// Define a custom testServer type which embeds an httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initializes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize the test server as normal.
	t_server := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client. Any response cookies will
	// now be stored in the jar and sent with the subsequent requests when using
	// this client.
	t_server.Client().Jar = jar

	// Prevent the test server client from following redirects by setting a
	// custom CheckRedirect function. This function runs whanever a 3xx
	// response is received. By returning http.ErrUseLastResponse, it tells
	// the client to stop and immediately return the received response.
	t_server.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{t_server}
}

func (t_server *testServer) resetClientCookieJar(t *testing.T) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	t_server.Client().Jar = jar

}

// Define a testResponse struct to hold data abaout responses from the test
// server. Note that this struct includes fields for the HTTP response headers
// and cookies, as well as the status code and body.
type testResp struct {
	status  int
	headers http.Header
	cookies []*http.Cookie
	body    string
}

// Implement a get() method on our custom testServer type. This makes a GET
// request to a given url path using the test server client and it returning a
// testResponse struct containing the response data.
func (t_server *testServer) get(t *testing.T, url_path string) testResp {
	req, err := http.NewRequest(http.MethodGet, t_server.URL+url_path, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := t_server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResp{res.StatusCode, res.Header, res.Cookies(), string(bytes.TrimSpace(body))}
}
