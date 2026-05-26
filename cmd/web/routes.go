package main

import "net/http"

// Update the signature for the routes() method so that it returns a 
// http.Handler instead of *http.ServeMux 

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	file_server := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", file_server))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Pass the servemux as the 'next' parameter to the commonHeaders middleware
	// Because commonHeaders is just a function, and the function returns a 
	// http.Handler we don´t need to do anything else.
	wrapped_handler := commonHeaders(mux) 

	return wrapped_handler
}
