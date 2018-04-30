package genie

import (
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NewLocalLambda creates a new lambda that executes off of local code
func NewLocalLambda(name, directory, command string, code []byte) (*LocalLambda, error) {
	l := &LocalLambda{
		name:    name,
		command: cmd(command) + " " + dir(directory) + name,
	}

	writeError := l.Write(directory+"/"+name, code)
	if writeError == nil {
		return l, nil
	}
	return l, writeError
}

// LocalLambda is a struct containing everything needed to write a file locally
type LocalLambda struct {
	name    string
	command string
}

// Name returns the string name of the filelambda
func (l *LocalLambda) Name() string {
	return l.name
}

// Write takes the code and writes it to the directory + name
func (l *LocalLambda) Write(file string, code []byte) error {
	if err := ioutil.WriteFile(file, code, 0755); err == nil {
		return nil
	}
	return errors.New("Unable to write file")
}

// Execute will execute the lambda and return it's output and errors(if applicable)
func (l *LocalLambda) Execute(stdin io.Reader, args []string) (string, error) {
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
		return string(stdoutStderr), nil
	}

	// *sigh*
	log.WithFields(log.Fields{"name": l.Name(), "command": l.command}).Error("Lambda Execution")
	return string(stdoutStderr), errors.New("Error running command")
}
