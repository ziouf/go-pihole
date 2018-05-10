package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"cm-cloud.fr/go-pihole/bdd"
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
	if v, ok := processMap[vars[`process`]]; ok {
		switch vars[`action`] {
		case `start`:
			v.Start()
		case `stop`:
			v.Stop()
		case `restart`:
			v.Restart()
		default:
		}
	}
}
