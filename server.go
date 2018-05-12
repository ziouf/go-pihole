package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"cm-cloud.fr/go-pihole/actions"
	"cm-cloud.fr/go-pihole/bdd"
	"cm-cloud.fr/go-pihole/config"
	"cm-cloud.fr/go-pihole/parser"
	"cm-cloud.fr/go-pihole/process"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

var (
	logTail *tail.Tail
	srv     *http.Server

	stop = make(chan os.Signal, 1)
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize application configurations
	// Load configuration files
	// Parse invocation flags
	// Parse environment variables
	config.Init()

	// Init embeded database
	bdd.Init()

	// Init managed processes
	process.Init()

	// Init HTTP Server
	srv = &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, initRouter()),
		Addr:    viper.GetString("bind"),
		// Good practice to enforce timeouts
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

func main() {
	defer close(stop)

	// Open DB
	bdd.Open()
	defer bdd.Close()

	// Start managed processes
	process.StartAll()
	defer process.StopAll()

	// Init dnsmasq log reader
	go logReaderService()

	// Start HTTP Server
	go httpServer()

	// Clean shutdown
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down the server...")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}

func initRouter() *mux.Router {
	// Router
	root := mux.NewRouter()

	// API route
	apiRoot := root.PathPrefix("/api/v1").Subrouter()

	// Model querying
	apiModel := apiRoot.PathPrefix("/model").Subrouter()

	// Get
	apiModel.HandleFunc(`/logs/dns/last`, lastDNSHandler)
	apiModel.HandleFunc(`/logs/dhcp/last`, lastDHCPHandler)
	// Search
	apiModel.HandleFunc(`/find/logs/dns/since/{date}`, nil)  /* FindLogsSinceDate */
	apiModel.HandleFunc(`/find/logs/dhcp/since/{date}`, nil) /* FindLogsSinceDate */
	// Stats
	apiModel.HandleFunc(`/stats/`, nil)

	// DHCP
	apiDhcp := apiRoot.PathPrefix("/dhcp").Subrouter()
	// DHCP Leases
	apiDhcp.HandleFunc("/leases", actions.AllLeases)                  /* List all dhcp leases */
	apiDhcp.HandleFunc("/leases/reserved", actions.AllReservedLeases) /* List all reserved dhcp leases */

	// Contrôler actions
	apiAction := apiRoot.PathPrefix("/action").Subrouter()
	// Actions on processes
	apiAction.HandleFunc("/process/{process}/{action:start|stop|restart}", processActionHandler) /*  */

	// Configuration
	apiConfig := apiRoot.PathPrefix(`/config`).Subrouter()
	// Update configuration values
	apiConfig.HandleFunc(``, nil)

	// Static content route
	root.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/dist/")))

	return root
}

func logReaderService() {
	file := viper.GetString("dnsmasq.log.file")
	if logTail, err := tail.TailFile(file, tail.Config{Follow: true, ReOpen: true}); err == nil {

		// Parallelise log parsing
		for i := 0; i < runtime.NumCPU()/2; i++ {
			go func() {
				for line := range logTail.Lines {
					bdd.Insert(parser.LogParser.ParseLine(line.Text))
				}
			}()
		}

	} else {
		log.Printf("Error while tailing file %s : %s", file, err)
	}
}

func httpServer() {
	log.Fatal(srv.ListenAndServe())
}
