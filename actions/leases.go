package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cm-cloud.fr/go-pihole/models"
	"cm-cloud.fr/go-pihole/utils"
	"github.com/spf13/viper"
)

// AllLeases return all DHCP leases
func AllLeases(w http.ResponseWriter, r *http.Request) {

	fileName := fmt.Sprintf("%s/%s", viper.GetString("dnsmasq-config-dir"), viper.GetString("dnsmasq-config-dhcp-leases"))

	lines, err := utils.ReadFileLines(fileName, func(line string) interface{} {
		lease := models.NewDhcpLease()
		lease.ParseLease(line)
		return lease
	})

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(lines)
}

// AllReservedLeases return all DHCP reserved leases from DNSMASQ config file
func AllReservedLeases(w http.ResponseWriter, r *http.Request) {
	// TODO : implement
	fileName := fmt.Sprintf("%s/%s", viper.GetString("dnsmasq-config-dir"), viper.GetString("dnsmasq-config-dhcp-static-leases"))

	lines, err := utils.ReadFileLines(fileName, func(line string) interface{} {
		lease := models.NewDhcpLease()
		lease.ParseStaticLease(line)
		return lease
	})

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(lines)
}
