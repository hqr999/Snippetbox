package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/hqr999/Snippetbox/internal/models"
)

// Define a snippetCreateForm struct to represent the form and validation 
// errors for the form fields. Note that all the struct fields are deliberately 
// exported (i.e start with a control letter). This is because struct fields
// must be exported in order to be read by the html/template package when 
// rendering the template. 
type snippetCreateForm struct {
		Title 					string 
		Content 				string 
		Expires 				int 
		FieldErrors map[string]string  
}



func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	//And the same thing again here
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data_tmpl := app.newTemplateData(r)
	
	// Initialize a new snippetCreateForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the 
	// snippet expiry to 365 days.
	data_tmpl.Form = snippetCreateForm{Expires: 365}

	app.render(w, r, http.StatusOK, "create.tmpl", data_tmpl)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	
	expiration_field_ops, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return

	}
	
	// Create an instance of the snippetCreateForm struct containing the values 
	// from the form and an empty map for any validation errors.
	form_data := snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expiration_field_ops,
		FieldErrors: map[string]string{},
}
	// Update the validation checks so that they operate on the snipperCreateForm 
	// instance.
	if strings.TrimSpace(form_data.Title) == "" {
		form_data.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form_data.Title) > 100 {
		form_data.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form_data.Content) == "" {
		form_data.FieldErrors["content"] = "This field can´t be blank"
	}

	if form_data.Expires != 1 && form_data.Expires != 7 && form_data.Expires != 365 {
		form_data.FieldErrors["expires"] = "This field must be equal to either 1, 7 or 365."
	}
	
	// If there are any validation errors, then the create.tmpl template, 
	// passing in the snipperCreateForm instance as dynamic data in the Form
	// field. Note  that we use the HTTP status code 422 Unprocessable Entity
	// when sending the response to indicate that there was a validation error.
	if len(form_data.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form_data
		app.render(w,r,http.StatusUnprocessableEntity,"create.tmpl",data)
		return
	}

	// We also need to update this line to pass the data from the 
	// snippetCreateForm instance to our Insert() method.
	id, err := app.snippets.Insert(form_data.Title,form_data.Content,form_data.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
