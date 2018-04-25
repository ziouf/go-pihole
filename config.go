package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// ApplicationName is the name of the application
const (
	ApplicationName = "go-pihole"
)

func init() {

	// Configure applicaiton
	configs()

}

func configs() {
	// Application
	viper.SetDefault("bind", ":8080")

	// SQlite
	viper.SetDefault("db_file", fmt.Sprintf("/var/lib/%1s/%1s.sqlite", ApplicationName))

	// DNMASQ
	viper.SetDefault("dnsmasq_pid_file", "/run/dnsmasq/dnsmasq.pid")
	viper.SetDefault("dnsmasq_resolv_file", "/run/dnsmasq/resolv.conf")
	viper.SetDefault("dnsmasq_log_file", "/var/log/dnsmasq.log")
	viper.SetDefault("dnsmasq_gravity_file", "/etc/pihole/gravity.list")
	viper.SetDefault("dnsmasq_config_dir", "/etc/dnsmasq.d")

	//

}
