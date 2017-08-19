package prox

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type Process interface {
	Name() string
	Run() error // TODO pass a ctx
	Interrupt() error
	RegisterEnvironment([]string)
}

type shellProcess struct {
	name   string
	script string
	env    []string
	logger *zap.Logger
	writer io.Writer

	mu  sync.Mutex
	cmd *exec.Cmd
}

func NewShellProcess(name, script string) Process {
	return &shellProcess{
		script: script,
		name:   name,
		logger: logger(name),
	}
}

func (p *shellProcess) Name() string {
	return p.name
}

func (p *shellProcess) RegisterEnvironment(env []string) {
	p.env = env
}

func (p *shellProcess) Run() error {
	p.mu.Lock()

	commandLine := p.buildCommandLine()
	p.logger.Debug("Starting process",
		zap.String("script", commandLine),
		zap.Strings("env", p.env),
	)

	cmdParts := strings.Split(commandLine, " ")
	p.cmd = exec.Command(cmdParts[0], cmdParts[1:]...)

	//writer := OutputWriter(output)
	//p.cmd.Stdout = writer
	//p.cmd.Stderr = writer
	p.cmd.Env = p.env

	err := p.cmd.Start()
	p.mu.Unlock()

	if err != nil {
		return fmt.Errorf("could not start shell task: %s", err)
	}

	return p.cmd.Wait()
}

func (p *shellProcess) buildCommandLine() string {
	return p.script

	//commandLine := p.Env.Replace(p.script)
	//
	//r := regexp.MustCompile(`[a-zA-Z_]+=\S+`)
	//
	//commandLineBuffer := new(bytes.Buffer)
	//parts := strings.Split(commandLine, " ") // TODO breaks if we have quotes spaces
	//done := false
	//for _, part := range parts {
	//	match := r.FindString(part)
	//	if done == false && match != "" {
	//		p.Env.Set(match)
	//	} else {
	//		done = true
	//	}
	//
	//	if done {
	//		commandLineBuffer.WriteString(part)
	//		commandLineBuffer.WriteString(" ")
	//	}
	//}
	//
	//return strings.TrimSpace(commandLineBuffer.String())
}

func (p *shellProcess) Interrupt() error {
	return errors.New("NOT IMPLEMENTED")
}
