package main

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
	return &Process{
		cmd: exec.Command(bin, args...),
	}
}

func (p *Process) Start() error {
	// Check if running
	if p.IsRunning() {
		return errors.New("Already running")
	}
	// If not running, start it
	if _, ok := processMap[p.cmd.Args[0]]; !ok {
		processMap[p.cmd.Args[0]] = p
	}
	if err := p.cmd.Run(); err != nil {
		return err
	}
	log.Printf("Starting [pid : %d] %s", p.cmd.Process.Pid, p.cmd.Args)
	return nil
}

func (p *Process) Stop() error {
	// Check if running
	if !p.IsRunning() {
		return errors.New("Not running")
	}
	// If running, stop it
	if _, ok := processMap[p.cmd.Args[0]]; !ok {
		return errors.New("Process not running")
	}
	return p.cmd.Process.Kill()
}

func (p *Process) Restart() error {
	if p.IsRunning() {
		p.Stop()
	}
	return p.Start()
}

func (p *Process) IsRunning() bool {
	return p.cmd.ProcessState.Success() && !p.cmd.ProcessState.Exited()
}

var processMap = make(map[string]*Process)

// ShutdownAll Stop all processes started by daemon
func ShutdownAll() {
	log.Println("Shuting down subprocesses")
	for _, value := range processMap {
		value.Stop()
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
