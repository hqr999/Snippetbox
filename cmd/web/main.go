package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr",":4000","HTTP network address")
	flag.Parse()
	
	// Use the slog.New() function to initialize a new structured logger, which
// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout,&slog.HandlerOptions{AddSource: true}))
		
	mux := http.NewServeMux()
	
	file_server := http.FileServer(http.Dir("./ui/static/"))

		mux.Handle("GET /static/", http.StripPrefix("/static",file_server))

	mux.HandleFunc("GET /{$}",home)
	mux.HandleFunc("GET /snippet/view/{id}",snippetView)
	mux.HandleFunc("GET /snippet/create",snippetCreate)
	mux.HandleFunc("POST /snippet/create",snippetCreatePost)
	
	
	logger.Info("starting server","addr",*addr)

	err := http.ListenAndServe(*addr,mux)
 	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServe() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}
