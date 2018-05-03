package process

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

var pp = make([]*Process, 0)

// Process Struct that represent a process
type Process struct {
	cmd *exec.Cmd
}

// NewProcess Create new Process
func NewProcess(bin string, args ...string) *Process {
	p := new(Process)
	p.cmd = exec.Command(bin, args...)
	return p
}

// Start the Process
func (p *Process) Start() error {
	// Check if running
	if p.isRunning() {
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

// Stop Kill the Process
func (p *Process) Stop() error {
	// Check if running
	if !p.isRunning() {
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

// Restart Restart the Process
func (p *Process) Restart() error {
	if p.isRunning() {
		p.Stop()
	}
	return p.Start()
}

func (p *Process) isRunning() bool {
	if p.cmd != nil && p.cmd.ProcessState != nil {
		return p.cmd.ProcessState.Success() && !p.cmd.ProcessState.Exited()
	}
	return false
}

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
