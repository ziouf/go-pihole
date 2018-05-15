package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"time"

	"cm-cloud.fr/go-pihole/log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Init default configuration values and configuration parsing
func Init() {
	setDefaults()
	parseEnv()
	parseFlags()
	parseConfigFile()

	// Init logger
	log.Init(path.Join(viper.GetString(`log.path`), viper.GetString(`log.file`)), viper.GetString(`log.level`))

	// Display configuration in debug logs
	log.Debug().Println(`== Config ==`)
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, k := range keys {
		v := viper.Get(k)
		t := reflect.TypeOf(v)
		switch t.Name() {
		case "int":
			log.Debug().Printf("%-40s : %d", k, v.(int))
		case "bool":
			log.Debug().Printf("%-40s : %t", k, v.(bool))
		default:
			log.Debug().Printf("%-40s : %s", k, v)
		}
	}
	log.Debug().Println(`============`)
}

func getApplicationPath() string {
	// Applications current path
	ex, err := os.Executable()
	if err != nil {
		log.Error().Fatalln(err)
	}
	p, err := filepath.EvalSymlinks(ex)
	if err != nil {
		log.Error().Fatalln(err)
	}
	app, err := filepath.Abs(filepath.Dir(p))
	if err != nil {
		log.Error().Fatalln(err)
	}
	return app
}

func setDefaults() {

	// Application
	viper.SetDefault(`app.name`, `go-pihole`)
	viper.SetDefault(`app.dir`, getApplicationPath())
	viper.SetDefault(`app.bind`, `:8080`)

	// Log
	viper.SetDefault(`log.level`, log.INFO)
	viper.SetDefault(`log.path`, "")
	viper.SetDefault(`log.file`, viper.GetString(`app.name`))

	// UI
	viper.SetDefault(`ui.path`, path.Join(viper.GetString(`app.dir`), `ui/dist/`))

	// Database
	viper.SetDefault(`db.file.path`, filepath.Join(viper.GetString(`app.dir`), fmt.Sprintf(`%s.db`, viper.GetString(`app.name`))))
	viper.SetDefault(`db.file.mode`, 0600)
	viper.SetDefault(`db.bulk.size`, 5000)
	viper.SetDefault(`db.bulk.freq`, time.Second)
	viper.SetDefault(`db.cleaning.enable`, false)
	viper.SetDefault(`db.cleaning.freq`, 250*time.Millisecond)
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
		// TODO : Add more
	})

	// DNSMASQ binary location
	viper.SetDefault(`dnsmasq.bin`, `/usr/sbin/dnsmasq`)
}

func parseFlags() {
	// Define flags
	flag.String(`log.level`, viper.GetString(`log.level`), `Logging level`)
	flag.String(`log.path`, viper.GetString(`log.path`), `Log output file. [Default : stderr]`)
	flag.String(`app.bind`, viper.GetString(`bind`), `IP:Port to bind HTTP server on`)
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
