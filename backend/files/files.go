package files

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"cm-cloud.fr/go-pihole/backend/log"
)

// Errors
var (
	ErrNotFoundInPath = errors.New(`Bin not found in $PATH`)
)

// ReadFileLines Read file and apply fn on each line
func ReadFileLines(fileName string, fn func(string) interface{}) ([]interface{}, error) {
	var err error

	result := make([]interface{}, 0)

	if f, err := os.Open(fileName); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			if line := fn(scanner.Text()); line != nil {
				result = append(result, line)
			}
		}
	}

	return result, err
}

// FindBinInPath Looks for binary present in PATH environment variable
// Return Full path of the binary or error
func FindBinInPath(bin string) (string, error) {
	for _, p := range strings.Split(os.Getenv("PATH"), ":") {
		path := fmt.Sprintf("%s/%s", p, bin)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else {
			log.Error().Println(err)
		}
	}
	return bin, ErrNotFoundInPath
}
