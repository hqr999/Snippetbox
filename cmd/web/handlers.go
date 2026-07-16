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
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, 200, "signup.tmpl", data)

}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare a zero-value instance of our userSignupForm struct
	var form userSignupForm

	// Parse the form data into the userSignupForm struct
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions
	form.CheckField(validator.NotBlank(form.Name), "name", "Field name cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Field email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must have a valid email like format")
	form.CheckField(validator.NotBlank(form.Password), "password", "Field password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Field password has to be at least 8 characters long")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "Field password cannot be more than 72 bytes long")

	// If there are any errors, redisplay the signup form along with a 422 status code
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and redisplay it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}

		return 
	}

	// Otherwise add a confirmation flash message to the session confirming that 
	// their signup worked. 
	app.sessionManager.Put(r.Context(),"flash","Your signup was successful. Please log in")

	// And redirect the user to the login page. 
	http.Redirect(w,r,"/user/login",http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
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

	err := app.decodePostForm(r, &form_data)
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
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
