package db

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"cm-cloud.fr/go-pihole/files/dnsmasq"
)

// LogsHandler Handle all logs
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	var result []dnsmasq.Log

	Db.Find(&result)

	json.NewEncoder(w).Encode(result)
}

// LogsLastHandler Handler last N logs
func LogsLastHandler(w http.ResponseWriter, r *http.Request) {
	var result []dnsmasq.Log
	vars := mux.Vars(r)

	Db.Limit(vars[`limit`]).Order("id desc").Find(&result)

	json.NewEncoder(w).Encode(result)
}

// LogSinceDateHandler Handle logs since date
func LogSinceDateHandler(w http.ResponseWriter, r *http.Request) {
	var result []dnsmasq.Log
	vars := mux.Vars(r)

	Db.Find(&result, "date >= ?", vars["date"])

	json.NewEncoder(w).Encode(result)
}

// LogHandler Handler one log
func LogHandler(w http.ResponseWriter, r *http.Request) {
	var result dnsmasq.Log
	vars := mux.Vars(r)

	Db.Find(&result, "id = ?", vars["id"])

	json.NewEncoder(w).Encode(result)
}
