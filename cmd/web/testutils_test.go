package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/hqr999/Snippetbox/internal/models/mocks"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	// Create an instance of the template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	// Add a form decoder
	formDecoder := form.NewDecoder()

	// And a session manager instance. Note that we use the same settings as
	// production, except that we don´t set a Store for session manager.
	// If no store is set, the SCS package will default to using a transient
	// in-memory store, which is ideal for testing purposes.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.DiscardHandler),
		snippets:       &mocks.SnippetModel{}, //Use the mock
		users:          &mocks.UserModel{},    //Use the mock
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
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

func extractCSRFToken(t *testing.T, body string) string {

	// Define a regex which captures the CSRF value from the
	// HTML for our signup page.
	csrfTokenRegex := regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

	// Use the FindStringSubmatch method to extract the token from the HTML body.
	// Note that this reutrns a slice with the entire matched pattern in the
	// first position, and the values of any captured data in the subsequent
	// positions.
	matches_found := csrfTokenRegex.FindStringSubmatch(body)
	if len(matches_found) < 2 {
		t.Fatal("no csrf on the html body.")

	}

	return html.UnescapeString(matches_found[1])
}

// Create a postForm method for sending POST requests to the test server. The
// final parameter to this method is a url.Values map which can contain any
// form data that you want to send in the request body.
func (test_server *testServer) postForm(t *testing.T, urlPath string, form url.Values) testResp {
	req, err := http.NewRequest(http.MethodPost, test_server.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	// Set the appropriate Content-Type header for form data and the Sec-Fetch-Site
	// header.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	res, err := test_server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResp{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}
