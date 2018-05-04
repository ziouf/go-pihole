package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cm-cloud.fr/go-pihole/actions"
	"cm-cloud.fr/go-pihole/bdd"
	"cm-cloud.fr/go-pihole/files/dnsmasq"
	"cm-cloud.fr/go-pihole/process"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

var (
	logTail        *tail.Tail
	dnsmasqProcess *process.Process
)

func init() {
	// Initialize application configurations
	// Load configuration files
	// Parse invocation flags
	// Parse environment variables
	initConfig()

	// Init embeded database
	bdd.Init()
	bdd.Open()
	bdd.AddToClean(&dnsmasq.Log{})
	// go bdd.CleanService(viper.GetDuration(`db.cleaning.freq`)*time.Second, viper.GetDuration(`db.cleaning.days.to.keep`))

	// Init dnsmasq log reader
	go logReaderService()

	// Init dnsmasq process
	if viper.GetBool(`dnsmasq.embeded`) {
		dnsmasqProcess = process.NewProcess(viper.GetString("dnsmasq.bin"),
			`-d`, `-k`, // No daemon
			`-C`, viper.GetString(`dnsmasq.config.file`),
			`-7`, fmt.Sprintf("%s,.dpkg-dist,.dpkg-old,.dpkg-new,.log,.sh,README", viper.GetString(`dnsmasq.config.dir`)),
			`-8`, viper.GetString(`dnsmasq.log.file`),
			// `-r`, fmt.Sprintf("%s/%s", viper.GetString(`dnsmasq.config.dir`), viper.GetString(`dnsmasq.config.resolv`)),
			// http://data.iana.org/root-anchors/root-anchors.xml
			`--trust-anchor=.,19036,8,2,49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5`,
			`--trust-anchor=.,20326,8,2,E06D44B80B8F1D39A95C0B0D7C65D08458E880409BBC683457104237C7F8EC8D`,
		)
		if err := dnsmasqProcess.Start(); err != nil {
			log.Fatalf("Error while starting DNSMASQ process : %s", err)
		}
	}

}

func main() {
	defer bdd.Close()
	defer process.ShutdownAll()

	// Router
	root := mux.NewRouter()

	// API route
	apiRoot := root.PathPrefix("/api/v1").Subrouter()

	// Model querying
	apiModel := apiRoot.PathPrefix("/model").Subrouter()
	apiModel.HandleFunc("/logs", nil)              /* AllLogs */
	apiModel.HandleFunc("/logs/last/{limit}", nil) /* AllLogs */
	apiModel.HandleFunc("/log/{id}", nil)          /* OneLog */
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

func logReaderService() {
	file := viper.GetString("dnsmasq.log.file")
	if logTail, err := tail.TailFile(file, tail.Config{Follow: true, ReOpen: true}); err == nil {

		for line := range logTail.Lines {
			bdd.Insert(dnsmasq.NewLog().ParseLine(line.Text))
		}

	} else {
		log.Printf("Error while tailing file %s : %s", file, err)
	}
}
