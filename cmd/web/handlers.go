package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)


// Change the signature of the home handler so it is defined as a method against
// *application.
func(app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	
	files_slice := []string{
				"./ui/html/base.tmpl",
				"./ui/html/pages/home.tmpl",
				"./ui/html/partials/nav.tmpl",
	}

	ts, err := template.ParseFiles(files_slice...)
	if err != nil {
			app.serverError(w,r,err) //Use the serverError() helper 
			return 
	}

	err = ts.ExecuteTemplate(w,"base",nil)
	if err != nil {
			app.serverError(w,r,err) //Use the serverError() helper		
	}

}

// Change the signature of the snippetView handler so it is defined as a method
// against *application.
func(app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}


// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func(app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}


// Change the signature of the snippetCreatePost handler so it is defined as a method
// against *application.
func(app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
