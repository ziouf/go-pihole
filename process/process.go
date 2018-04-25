package process

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"cm-cloud.fr/go-pihole/utils"
)

//
const (
	Start int8 = iota
	Stop
	Restart
)

var processMap = make(map[string]*exec.Cmd, 0)

// ShutdownAll Stop all processes started by daemon
func ShutdownAll() {
	log.Println("Shuting down subprocesses")
	for key := range processMap {
		Process(key, Stop)
	}
	log.Println("Subprocesses stopped")
}

// Process Start/Stop/Restart managed Process
func Process(name string, action int8, args ...string) error {
	bin, err := utils.FindBinInPath(name)
	if err != nil {
		return fmt.Errorf("Can't Start/Stop process because: %s", err)
	}

	switch action {
	case Restart:
		err = Process(name, Stop, args...)
		err = Process(name, Start, args...)

	case Start:
		sh, ok := processMap[bin]
		if !ok {
			if b, pid := IsProcessRunning(name); b {
				log.Printf("%s is running [pid : %d] -> Killing it", name, pid)
				if p, err := os.FindProcess(pid); err == nil {
					if err := p.Kill(); err != nil {
						log.Fatal(err)
					}
				}
			}
			processMap[bin] = exec.Command(bin, args...)
			sh = processMap[bin]
		}
		if err := sh.Run(); err != nil {
			log.Printf("Error : %s", err)
		}
		log.Printf("Starting [pid : %d] %s", sh.Process.Pid, sh.Args)

	case Stop:
		if sh, ok := processMap[bin]; ok {
			log.Printf("Killing process '%s' [pid : %d]", sh.Path, sh.Process.Pid)
			sh.Process.Kill()
		}
	}

	return err
}

// IsProcessRunning Check if there is at least one process of the specified program running
func IsProcessRunning(name string) (bool, int) {
	var running, pid, root = false, 0, "/proc"

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return fmt.Errorf("%s is not a Directory", path)
		}
		if match, err := regexp.MatchString("^[0-9]+$", info.Name()); err == nil && match {
			// Read file comm in subdir
			utils.ReadFileLines(fmt.Sprintf("%s/comm", path), func(line string) interface{} {
				if line == name {
					running = true
					i, _ := strconv.ParseInt(strings.Replace(path, fmt.Sprintf("%s/", root), "", 1), 10, 64)
					pid = int(i)
				}
				return nil
			})
			return filepath.SkipDir
		}
		return nil
	})
	return running, pid
}
