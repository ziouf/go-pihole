package files

import (
	"strings"
)

// DNSMASQ config mapping
type config struct {
}

// DNSMASQ configuration files parsing and writing
func readConfigFile(fname string) {
	// get lines filtering comments and blank
	lines, err := ReadFileLines(fname, func(line string) interface{} {
		if strings.HasPrefix(`#`, line) {
			return nil
		}
		if len(strings.TrimSpace(line)) == 0 {
			return nil
		}
		return line
	})

	// manage file reader error
	if err != nil {

	}

	// get configuration key/values
	config := make(map[string]string, 0)
	for _, line := range lines {
		tokens := strings.Split(line.(string), `=`)
		config[tokens[0]] = tokens[1]
	}

	// save configuration key/values to datastore
	//...

}
