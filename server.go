package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	// Router
	root := mux.NewRouter()

	// API route
	// apiRoot := root.PathPrefix("/api/v1/").Subrouter()
	// apiRoot.HandleFunc("/logs", api.LogsHandler)

	// Static content route
	root.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/dist/")))

	// Http server
	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, root),
		Addr:    viper.GetString("bind"),
		// Good practice to enforce timeouts
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
