package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hqr999/Snippetbox/internal/models"
	"github.com/hqr999/Snippetbox/internal/validator"
)


type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	validator.Validator `form:"-"`
}


func (app *application) userSignup(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"Display a form for signing up a new user...")
}

func (app *application) userSignupPost(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"Display a form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter,r *http.Request)  {
	fmt.Fprintln(w,"Logout the user...")
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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data_tmpl := app.newTemplateData(r)

	data_tmpl.Form = snippetCreateForm{Expires: 365}

	app.render(w, r, http.StatusOK, "create.tmpl", data_tmpl)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form_data snippetCreateForm

	err := app.decodePostForm(r,&form_data)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form_data.CheckField(validator.NotBlank(form_data.Title), "title", "This field cannot be blank")
	form_data.CheckField(validator.MaxChars(form_data.Title, 100), "title", "This field cannot be bigger than 10 characters")
	form_data.CheckField(validator.NotBlank(form_data.Content), "content", "This field cannot be blank")
	form_data.CheckField(validator.PermittedValue(form_data.Expires, 1, 7, 365), "expires", "This field must be equal 1, 7 or 365")

	if !form_data.Valid() {
		data := app.newTemplateData(r)
		data.Form = form_data
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form_data.Title, form_data.Content, form_data.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//Use the Put() method to add string value ("Snippet successfully"
	//created") and the corresponding key ("flash") to the session data.
	app.sessionMangaer.Put(r.Context(),"flash","Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
