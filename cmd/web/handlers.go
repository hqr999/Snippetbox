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

	app.render(w, r, http.StatusOK, "create.tmpl", data_tmpl)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title_field := r.PostForm.Get("title")
	content_field := r.PostForm.Get("content")

	expiration_field_ops, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return

	}
	// Initialize a map to hold any validation errors fir the form fields
	formFieldErrors := make(map[string]string)

	// Check that the title value is not blank and is
	// not more than 100 characters of lenght. If it
	// falls either of those checks, add a message
	// to the errors map using the field name as the key.
	if strings.TrimSpace(title_field) == "" {
		formFieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title_field) > 100 {
		formFieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check that the content value isn´t blank
	if strings.TrimSpace(content_field) == "" {
		formFieldErrors["content"] = "This field can´t be blank"
	}

	// Check the expires value matches of the permitted values
	// (1, 7, or 365)
	if expiration_field_ops != 1 && expiration_field_ops != 7 && expiration_field_ops != 365 {
		formFieldErrors["expiration option"] = "This field must be equal to either 1, 7 or 365."
	}

	// If there are any errors, dump them in a plain.Text HTTP response and
	// return from the handler.
	if len(formFieldErrors) > 0 {
		fmt.Fprint(w, formFieldErrors)
		return
	}

	id, err := app.snippets.Insert(title_field, content_field, expiration_field_ops)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
