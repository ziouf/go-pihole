package process

import (
	"fmt"

	"github.com/spf13/viper"
)

type Key string

const (
	// DNSMASQ process key
	DNSMASQ = Key(`dnsmasq`)
)

var pMap = make(map[Key]*Process)

// Init all processes
// Keep in mind that it does not start processes
func Init() {
	// DNSMASQ process int
	if viper.GetBool(`dnsmasq.embeded`) {
		pMap[DNSMASQ] = NewProcess(viper.GetString("dnsmasq.bin"),
			`-d`, `-k`, // No daemon
			`-C`, viper.GetString(`dnsmasq.config.file`),
			`-7`, fmt.Sprintf("%s,.dpkg-dist,.dpkg-old,.dpkg-new,.log,.sh,README", viper.GetString(`dnsmasq.config.dir`)),
			`-8`, viper.GetString(`dnsmasq.log.file`),
			// `-r`, fmt.Sprintf("%s/%s", viper.GetString(`dnsmasq.config.dir`), viper.GetString(`dnsmasq.config.resolv`)),
			// http://data.iana.org/root-anchors/root-anchors.xml
			`--trust-anchor=.,19036,8,2,49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5`,
			`--trust-anchor=.,20326,8,2,E06D44B80B8F1D39A95C0B0D7C65D08458E880409BBC683457104237C7F8EC8D`,
		)
	}
	// Other processes below
	// ...
}

func Start(name Key) error {
	if v, ok := pMap[name]; !ok {
		return fmt.Errorf(`Process [%s] not found`, name)
	} else {
		return v.start()
	}
}

func Restart(name Key) error {
	if v, ok := pMap[name]; !ok {
		return fmt.Errorf(`Process [%s] not found`, name)
	} else {
		return v.restart()
	}
}

func Stop(name Key) error {
	if v, ok := pMap[name]; !ok {
		return fmt.Errorf(`Process [%s] not found`, name)
	} else {
		return v.stop()
	}
}

func StartAll() {
	for _, v := range pMap {
		v.start()
	}
}

func RestartAll() {
	for _, v := range pMap {
		v.restart()
	}
}

func StopAll() {
	for _, v := range pMap {
		v.stop()
	}
}
