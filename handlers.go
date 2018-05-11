package main

import (
	"encoding/json"
	"log"
	"net/http"

	"cm-cloud.fr/go-pihole/bdd"
	"cm-cloud.fr/go-pihole/process"
	"github.com/gorilla/mux"
)

// LastDNSHandler Finds the most recent DNS log entry
func LastDNSHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := bdd.GetLast(new(bdd.DNS))
		if err != nil {
			log.Println(err)
		}
		json.NewEncoder(w).Encode(data)
	default:
	}
}

// LastDHCPHandler Finds the most recent DHCP log entry
func LastDHCPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := bdd.GetLast(new(bdd.DHCP))
		if err != nil {
			log.Println(err)
		}
		json.NewEncoder(w).Encode(data)
	default:
	}
}

func processActionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch vars[`action`] {
	case `start`:
		process.Start(process.Key(vars[`process`]))
	case `restart`:
		process.Restart(process.Key(vars[`process`]))
	case `stop`:
		process.Stop(process.Key(vars[`process`]))
	default:
	}
}
