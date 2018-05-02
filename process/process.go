package process

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"cm-cloud.fr/go-pihole/files"
)

type Process struct {
	cmd *exec.Cmd
}

func NewProcess(bin string, args ...string) *Process {
	p := new(Process)
	p.cmd = exec.Command(bin, args...)
	return p
}

func (p *Process) Start() error {
	// Check if running
	if p.IsRunning() {
		return errors.New("Already running")
	}
	// If not running, start it
	if err := p.cmd.Start(); err != nil {
		return err
	}
	pp = append(pp, p)
	log.Printf("Starting [pid : %d] %s", p.cmd.Process.Pid, p.cmd.Args)
	return nil
}

func (p *Process) Stop() error {
	// Check if running
	if !p.IsRunning() {
		return fmt.Errorf("Not running [pid : %d]", p.cmd.Process.Pid)
	}

	if err := p.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("Can't kill process [pid: %d] : %s", p.cmd.Process.Pid, err)
	}

	for i, v := range pp {
		if v == p {
			pp = append(pp[:i], pp[i+1:]...)
			break
		}
	}
	return nil
}

func (p *Process) Restart() error {
	if p.IsRunning() {
		p.Stop()
	}
	return p.Start()
}

func (p *Process) IsRunning() bool {
	if p.cmd != nil && p.cmd.ProcessState != nil {
		return p.cmd.ProcessState.Success() && !p.cmd.ProcessState.Exited()
	}
	return false
}

var pp = make([]*Process, 0)

// ShutdownAll Stop all processes started by daemon
func ShutdownAll() {
	log.Println("Shuting down subprocesses")

	for _, p := range pp {
		if err := p.Stop(); err != nil {
			log.Println(err)
		}
	}

	log.Println("Subprocesses stopped")
}

// IsProcessRunning Check if there is at least one process of the specified program running
func isProcessRunning(name string) (bool, int) {
	var running, pid, root = false, 0, "/proc"

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return fmt.Errorf("%s is not a Directory", path)
		}
		if match, err := regexp.MatchString("^[0-9]+$", info.Name()); err == nil && match {
			// Read file comm in subdir
			files.ReadFileLines(fmt.Sprintf("%s/comm", path), func(line string) interface{} {
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
