package main 


import "github.com/hqr999/Snippetbox/internal/models"


// Include a Snippets field in the templateData struct.
type templateData struct{
		Snippet models.Snippet
		Snippets []models.Snippet
}
