package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)


func snippetCreatePost(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Save a new snippet..."))
}

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Add("Server","Go")

	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	//w.Write([]byte(msg)) old way 
	fmt.Fprintf(w,"Display a specific snippet with ID %d...", id)

}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux_router := http.NewServeMux()
	// Think of the dollar tree as a wildcard *
	//Without it anything you type on the route will be
	//redirected to go with home function
	mux_router.HandleFunc("GET /{$}", home)
	mux_router.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux_router.HandleFunc("GET /snippet/create", snippetCreate)
	mux_router.HandleFunc("POST /snippet/create",snippetCreatePost)

	// Print a log message to say that the server is starting.
	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux_router)
	log.Fatal(err)
}
