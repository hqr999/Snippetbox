package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/hqr999/Snippetbox/internal/models"
)

// Include a Snippets field in the templateData struct.
type templateData struct {
	CurrYear int
	Snippet  models.Snippet
	Snippets []models.Snippet
	Form 		 any
	Flash    string //Add a Flash field to the templateData struct.
}


// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

//Initialize a template.FuncMap object and store it in  a global variable. This is 
//essentially a string-keyed map which acts as a lookup between the names of our 
//custom template functions and the functions themselves.
var funcs = template.FuncMap{
		"human_date": humanDate,
}


func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pgs, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, val := range pgs {
		name := filepath.Base(val)
		
		// The template.FuncMap must be registered with the template set before you
   // call the ParseFiles() method. This means we have to use template.New() to
  // create an empty template set, use the Funcs() method to register the
 // template.FuncMap, and then parse the file as normal.
		temp_s, err := template.New(name).Funcs(funcs).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		temp_s, err = temp_s.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		temp_s, err = temp_s.ParseFiles(val)

		cache[name] = temp_s

	}

	return cache, nil

}
