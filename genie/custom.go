package genie

import (
	"errors"
	"io"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NewCustomLambda creates a new lambda that has a custom command
func NewCustomLambda(name, command string) *CustomLambda {
	return &CustomLambda{
		name:    name,
		command: command,
	}
}

// CustomLambda creates a custom lambda(a command essentially)
type CustomLambda struct {
	name    string
	command string
}

// Name returns the custom lambda's name
func (l *CustomLambda) Name() string {
	return l.name
}

// Execute executes the custom lambda command
func (l *CustomLambda) Execute(stdin io.Reader, args []string) (string, error) {
	argsStr := strings.TrimSpace(strings.Join(args, " "))
	if argsStr != "" {
		argsStr = " " + argsStr
	}

	cmd := exec.Command("bash", "-c", l.command+argsStr)

	// pass through some stdin goodness
	cmd.Stdin = stdin

	// for those who are about to rock, I salute you.
	stdoutStderr, err := cmd.CombinedOutput()

	if err == nil {
		// noiiiice!
		log.WithFields(log.Fields{"name": l.Name(), "command": l.command}).Info("Lambda Execution")
		return strings.TrimSpace(string(stdoutStderr)), nil
	}

	// *sigh*
	log.WithFields(log.Fields{"name": l.Name(), "command": l.command}).Error("Lambda Execution")
	return string(stdoutStderr), errors.New("Error running command")
}
