package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init() {
	setDefaults()
	parseEnv()
	parseFlags()
	parseConfigFile()
}

func getApplicationPath() string {
	// Applications current path
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	p, err := filepath.EvalSymlinks(ex)
	if err != nil {
		log.Fatalln(err)
	}
	app, err := filepath.Abs(filepath.Dir(p))
	if err != nil {
		log.Fatalln(err)
	}
	return app
}

func setDefaults() {
	// Application
	viper.SetDefault(`app.name`, `go-pihole`)
	viper.SetDefault(`app.dir`, getApplicationPath())
	viper.SetDefault(`bind`, `:8080`)

	// SQlite
	viper.SetDefault(`db.file.path`, filepath.Join(viper.GetString(`app.dir`), fmt.Sprintf(`%s.db`, viper.GetString(`app.name`))))
	viper.SetDefault(`db.file.mode`, 0600)
	viper.SetDefault(`db.bulk.size`, 5000)
	viper.SetDefault(`db.bulk.freq`, 500*time.Millisecond)
	viper.SetDefault(`db.cleaning.enable`, false)
	viper.SetDefault(`db.cleaning.freq`, 10*time.Millisecond)
	viper.SetDefault(`db.cleaning.keep`, 7*time.Hour*24)

	// DNMASQ
	viper.SetDefault(`dnsmasq.embeded`, false)
	viper.SetDefault(`dnsmasq.pid.file`, `/run/dnsmasq/dnsmasq.pid`)
	viper.SetDefault(`dnsmasq.log.file`, `/var/log/dnsmasq.log`)

	viper.SetDefault(`dnsmasq.config.file`, `/etc/dnsmasq.conf`)
	viper.SetDefault(`dnsmasq.config.dir`, `/etc/dnsmasq.d`)

	viper.SetDefault(`dnsmasq.config.files.app`, `00_app.conf`)
	viper.SetDefault(`dnsmasq.config.files.resolv`, `01_resolv.conf`)
	viper.SetDefault(`dnsmasq.config.files.gravity`, `05_gravity.conf`)
	viper.SetDefault(`dnsmasq.config.files.dns`, `10_dns_cache.conf`)
	viper.SetDefault(`dnsmasq.config.files.dhcp-global`, `20_dhcp.conf`)
	viper.SetDefault(`dnsmasq.config.files.dhcp-static-leases`, `21_dhcp_static_leases.conf`)
	viper.SetDefault(`dnsmasq.config.files.dhcp-dynamic-leases`, `22_dhcp_leases.conf`)

	viper.SetDefault(`dnsmasq.gravity.lists`, []string{
		`https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts`,
		`https://mirror1.malwaredomains.com/files/justdomains`,
		// ...
	})

	// DNSMASQ binary location
	viper.SetDefault(`dnsmasq.bin`, `/usr/sbin/dnsmasq`)
}

func parseFlags() {
	// Define flags
	flag.String(`bind`, viper.GetString(`bind`), `IP:Port to bind HTTP server on`)
	flag.Bool(`db.cleaning.enable`, viper.GetBool(`db.cleaning.enable`), `Enable database auto cleaning`)
	flag.String(`dnsmasq.log.file`, viper.GetString(`dnsmasq.log.file`), `Dnsmasq log file`)
	flag.String(`dnsmasq.config.dir`, viper.GetString(`dnsmasq.config.dir`), `Dnsmasq configuration files directory`)
	flag.String(`dnsmasq.config.file`, viper.GetString(`dnsmasq.config.file`), `Dnsmasq configuration files directory`)
	flag.Bool(`dnsmasq.embeded`, viper.GetBool(`dnsmasq.embeded`), `Use embeded dnsmasq`)

	// Parse flags
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func parseEnv() {
	viper.SetEnvPrefix(`pihole`)
	viper.AutomaticEnv()
}

func parseConfigFile() {
	// Config file name
	viper.SetConfigName(`config`)

	viper.AddConfigPath(viper.GetString(`app.dir`))
	viper.AddConfigPath(filepath.Join(`etc`, fmt.Sprintf(`%[1]s`, viper.GetString(`application.name`))))
	viper.AddConfigPath(filepath.Join(`$HOME`, fmt.Sprintf(`.%[1]s`, viper.GetString(`application.name`))))

	if err := viper.ReadInConfig(); err == nil {
		viper.WatchConfig()
	}
}
