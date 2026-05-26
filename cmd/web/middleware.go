package main

import "net/http"

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
