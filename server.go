package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cm-cloud.fr/go-pihole/actions"
	"cm-cloud.fr/go-pihole/configs"
	"cm-cloud.fr/go-pihole/models"
	"cm-cloud.fr/go-pihole/process"
	"cm-cloud.fr/go-pihole/services"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func init() {
	configs.InitConfig()
	models.InitDB()

	services.StartLogReaderService()
}

func main() {
	defer models.Db.Close()
	defer process.ShutdownAll()

	// Router
	root := mux.NewRouter()

	// API route
	apiRoot := root.PathPrefix("/api/v1").Subrouter()

	// Model querying
	apiModel := apiRoot.PathPrefix("/model").Subrouter()
	apiModel.HandleFunc("/logs", nil)     /* AllLogs */
	apiModel.HandleFunc("/log/{id}", nil) /* OneLog */
	// Search
	apiModel.HandleFunc("/find/logs/since/{date}", nil)          /* FindLogsSinceDate */
	apiModel.HandleFunc("/find/logs/since/hour", nil)            /* FindLogsSinceAnHour */
	apiModel.HandleFunc("/find/logs/since/hour/{hours}", nil)    /* FindLogsSinceHours */
	apiModel.HandleFunc("/find/logs/since/day", nil)             /* FindLogsSinceADay */
	apiModel.HandleFunc("/find/logs/since/day/{days}", nil)      /* FindLogsSinceDays */
	apiModel.HandleFunc("/find/logs/between/{start}/{end}", nil) /* FindLogsBetweenDates */
	// Stats
	apiModel.HandleFunc("/stats/logs/count/type/since/{date}", nil)            /*  */
	apiModel.HandleFunc("/stats/logs/count/type/between/{start}/{end}", nil)   /*  */
	apiModel.HandleFunc("/stats/logs/count/client/since/{date}", nil)          /*  */
	apiModel.HandleFunc("/stats/logs/count/client/between/{start}/{end}", nil) /*  */
	apiModel.HandleFunc("/stats/logs/count/qtype/since/{date}", nil)           /*  */
	apiModel.HandleFunc("/stats/logs/count/qtype/between/{start}/{end}", nil)  /*  */

	// DHCP
	apiDhcp := apiRoot.PathPrefix("/dhcp").Subrouter()
	// DHCP Leases
	apiDhcp.HandleFunc("/leases", actions.AllLeases)                  /* List all dhcp leases */
	apiDhcp.HandleFunc("/leases/reserved", actions.AllReservedLeases) /* List all reserved dhcp leases */

	// Contrôler actions
	apiAction := apiRoot.PathPrefix("/action").Subrouter()
	apiAction.HandleFunc("/dnsmasq/restart", nil) /*  */

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

	go func() {
		log.Println("Server started...")
		log.Fatal(srv.ListenAndServe())
	}()

	stop := make(chan os.Signal, 1)
	defer close(stop)

	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}
