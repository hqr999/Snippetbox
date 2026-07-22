package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf" // New import
)

// Create a preventCSRF middleware function which uses a custom CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func preventCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func commonHeaders(prx http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is split across multiple lines for readability
		//purposes. Its not an obligation.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nonsniff")
		w.Header().Set("X-Frama-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		prx.ServeHTTP(w, r)
	})

}

func (dx *application) logRequest(prx http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		dx.logger.Info("request has been received", "ip", ip, "proto", proto, "method", method, "uri", uri)

		prx.ServeHTTP(w, r)
	})

}

func (app *application) recoverPanic(prx http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be r
		//	un in the event of a panic as Go unwinds the stack)
		defer func() {
			// Use the built-in recover func to check if there has
			// been a panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response;
				w.Header().Set("Connection", "close")
				//Cal the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, r, fmt.Errorf("%s", err))

			}

		}()
		prx.ServeHTTP(w, r)

	})

}

func (app *application) requireAuthentication(nxt http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authorized, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache ( or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain.
		nxt.ServeHTTP(w, r)

	})
}

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticateUserID value from the session using the
		// GetInt() method. This will return thr zero value for an int (0) if no
		// "authenticatedUserID" value is in the session - in which case we
		// call the next handler in the chain as normal and return.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we check to see if a user with that ID exists in our
		// database.
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// If a matching user is found, we know that the requestr is
		// coming from an authenticated user who exists in our database. We
		// create a new copy of the request (with an isAuthenticatedContextKey
		// value of true in the request context) and assign it to r.
		if exists {
			contx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(contx)
		}

		// Call the next handler
		next.ServeHTTP(w, r)

	})

}
