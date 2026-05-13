package main 


import "github.com/hqr999/Snippetbox/internal/models"


// Define a templateData type to act as the holding strcuture for  
// any dynamic data that we want to pass to our HTML templates.
//At the moment it only contains one field, but we will add more 
//to it as the build goes on.
type templateData struct{
		Snippet models.Snippet
}
