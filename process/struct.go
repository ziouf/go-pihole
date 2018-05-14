package process

import (
	"errors"
	"os/exec"

	"cm-cloud.fr/go-pihole/log"
)

// Errors
var (
	ErrRunning      = errors.New(`Already running`)
	ErrNotRunning   = errors.New(`Not running`)
	ErrCantKillProc = errors.New(`Can't kill process`)
)

// Process Struct that represent a process
type Process struct {
	cmd *exec.Cmd
}

// NewProcess Create new Process
func NewProcess(bin string, args ...string) *Process {
	log.Debug().Printf("Creating process %s with args %s", bin, args)
	p := new(Process)
	p.cmd = exec.Command(bin, args...)
	return p
}

// Start the Process
func (p *Process) start() error {
	log.Debug().Printf("Starting %s", p.cmd.Args)
	// Check if running
	if p.isRunning() {
		return ErrRunning
	}
	// If not running, start it
	if err := p.cmd.Start(); err != nil {
		return err
	}
	log.Info().Printf("Starting [pid : %d] %s", p.cmd.Process.Pid, p.cmd.Args)
	return nil
}

// Stop Kill the Process
func (p *Process) stop() error {
	log.Debug().Printf("Stopping %s", p.cmd.Args)
	// Check if running
	if !p.isRunning() {
		return ErrNotRunning
	}

	if err := p.cmd.Process.Kill(); err != nil {
		log.Error().Printf("Can't kill process [pid: %d] : %s", p.cmd.Process.Pid, err)
		return ErrCantKillProc
	}
	return nil
}

// Restart Restart the Process
func (p *Process) restart() error {
	log.Debug().Printf("Restarting %s", p.cmd.Args)
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
