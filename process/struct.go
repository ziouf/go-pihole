package process

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

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
func (p *Process) start() error {
	// Check if running
	if p.isRunning() {
		return errors.New("Already running")
	}
	// If not running, start it
	if err := p.cmd.Start(); err != nil {
		return err
	}
	log.Printf("Starting [pid : %d] %s", p.cmd.Process.Pid, p.cmd.Args)
	return nil
}

// Stop Kill the Process
func (p *Process) stop() error {
	// Check if running
	if !p.isRunning() {
		return fmt.Errorf("Not running [pid : %d]", p.cmd.Process.Pid)
	}

	if err := p.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("Can't kill process [pid: %d] : %s", p.cmd.Process.Pid, err)
	}
	return nil
}

// Restart Restart the Process
func (p *Process) restart() error {
	if p.isRunning() {
		p.stop()
	}
	return p.start()
}

func (p *Process) isRunning() bool {
	if p.cmd != nil && p.cmd.ProcessState != nil {
		return p.cmd.ProcessState.Success() && !p.cmd.ProcessState.Exited()
	}
	return false
}
