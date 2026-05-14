package main

import (
	"html/template"
	"path/filepath"

	"github.com/hqr999/Snippetbox/internal/models"
)

// Include a Snippets field in the templateData struct.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	//Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	//Use the filepath.Glob() function to get a slice of all filepaths that
	//match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	//us a slice of all th filepaths for our application 'page' templates
	//like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pgs, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, val := range pgs {
		//Extract the file name (like 'home.tmpl') from the full filepath
		//and assign it to the name variable.
		name := filepath.Base(val)

		//Parse the base template file into a template set.
		temp_s, err := template.ParseFiles("./ui/home/base.tmpl")
		if err != nil {
			return nil, err
		}

		//Call the parseGlob() * on this template set to add any partials.
		temp_s, err = temp_s.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		//Call ParseFiles() *on this template set* to add the page template.
		temp_s, err = temp_s.ParseFiles(val)

		//Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key
		cache[name] = temp_s

	}

	//Return the map
	return cache, nil

}
