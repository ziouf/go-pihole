package main

import (
	"encoding/json"
	"log"
	"net/http"

	"cm-cloud.fr/go-pihole/bdd"
	"cm-cloud.fr/go-pihole/process"
	"github.com/gorilla/mux"
)

func lastDNSHandler(w http.ResponseWriter, r *http.Request) {
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

func lastDHCPHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Println(process.Start(process.Key(vars[`process`])))
	case `restart`:
		log.Println(process.Restart(process.Key(vars[`process`])))
	case `stop`:
		log.Println(process.Stop(process.Key(vars[`process`])))
	default:
	}
}


func statCountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch vars[`type`] {
	case `dns`:
		if c, err := bdd.Count(&bdd.DNS{}); err != nil {
			log.Println(err)
		} else {
			json.NewEncoder(w).Encode(c)
		}
	case `dhcp`:
		if c, err := bdd.Count(&bdd.DHCP{}); err != nil {
			log.Println(err)
		} else {
			json.NewEncoder(w).Encode(c)
		}
	default:

	}
}