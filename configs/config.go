package configs

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// InitConfig ...
func InitConfig() {
	setDefaults()
	parseEnv()
	parseFlags()
	parseConfigFile()
}

func setDefaults() {
	// Application
	viper.SetDefault("application-name", "go-pihole")
	viper.SetDefault("bind", ":8080")

	// SQlite
	// viper.SetDefault("db_file", fmt.Sprintf("/var/lib/%[1]s/%[1]s.sqlite", viper.GetString("application-name")))
	viper.SetDefault("db_file", fmt.Sprintf("./%[1]s.sqlite", viper.GetString("application-name")))

	// DNMASQ
	viper.SetDefault("dnsmasq-pid-file", "/run/dnsmasq/dnsmasq.pid")
	viper.SetDefault("dnsmasq-log-file", "/var/log/dnsmasq.log")
	// viper.SetDefault("dnsmasq-gravity-file", "/etc/pihole/gravity.list")
	// viper.SetDefault("dnsmasq-dhcp-lease-file", "/etc/pihole/dhcp.leases")

	viper.SetDefault("dnsmasq-config-dir", "/etc/dnsmasq.d")
	viper.SetDefault("dnsmasq-config-app", "00_app.conf")
	viper.SetDefault("dnsmasq-config-resolv", "01_resolv.conf")
	viper.SetDefault("dnsmasq-config-gravity", "05_gravity.conf")
	viper.SetDefault("dnsmasq-config-dns", "10_dns_cache.conf")
	viper.SetDefault("dnsmasq-config-dhcp", "20_dhcp.conf")
	viper.SetDefault("dnsmasq-config-dhcp-static-leases", "21_dhcp_static_leases.conf")
	viper.SetDefault("dnsmasq-config-dhcp-leases", "22_dhcp_leases.conf")

	// DNSMASQ pocess arguments
	viper.SetDefault("dnsmasq-process-args", []string{
		"-x", viper.GetString("dnsmasq-pid-file"),
		// "-u", "dnsmasq",
		"-r", fmt.Sprintf("%s/%s", viper.GetString("dnsmasq-config-dir"), viper.GetString("dnsmasq-config-resolv")),
		"-7", fmt.Sprintf("%s,.dpkg-dist,.dpkg-old,.dpkg-new", viper.GetString("dnsmasq-config-dir")),
		"--local-service",
		"--trust-anchor=.,19036,8,2,49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5",
	})
}

func parseFlags() {
	// Define flags
	flag.String("bind", viper.GetString("bind"), "IP:Port to bind HTTP server on")
	flag.String("dhcp-leases-file", viper.GetString("dnsmasq-dhcp-lease-file"), "DHCP leases file")
	flag.String("dnsmasq-log-file", viper.GetString("dnsmasq-log-file"), "Dnsmasq log file")
	flag.String("dnsmasq-config-dir", viper.GetString("dnsmasq-config-dir"), "Dnsmasq configuration files directory")

	// Parse flags
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func parseEnv() {
	viper.SetEnvPrefix("pihole")
	viper.AutomaticEnv()
}

func parseConfigFile() {
	// Config file name
	viper.SetConfigName("config")

	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("/etc/%[1]s", viper.GetString("application-name")))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%[1]s", viper.GetString("application-name")))

	if err := viper.ReadInConfig(); err == nil {
		viper.WatchConfig()
	}
}
