package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "url", url, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)

}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrYear: time.Now().Year(),
		//Use the PopString() method to retrieve the value for the "flash" key.
		//PopString() also deletes the key and value from the session data, so it
		//behaves like a one-time fetch. If there is no matching key in the session
		//data this will return the empty string.

		//Add the flash message to the template data, if one exists
		Flash: app.sessionMangaer.PopString(r.Context(), "flash"),
	}

}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination into which we want to decode the form data.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	//Call ParseForm() on the request, in the same way that we did in our
	// SnippetCreatePost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	//Call Decode() on our decoder instance, pass the target destination as
	//the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		//If we try to use an invalid target destination, the Decode() method
		//will return an error with the type form.InvalidDcodeError. We use
		//errors.AsType() to check for this and panic. At the end of this
		//chapter we will talk about panicking versus returning errors, and
		//discuss why it is an appropriate thing to do this specific situation.
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}
	}

	//For all other errors, return them as normal
	return err
}
