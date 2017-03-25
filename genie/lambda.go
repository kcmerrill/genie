package genie

import (
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NewLambda creates a new lambda based off of text
func NewLambda(dir, name, cmd string, code []byte) (*Lambda, error) {
	l := &Lambda{
		Name:    name,
		Code:    code,
		Dir:     dir,
		Command: cmd,
	}

	writeError := l.Write()
	if writeError == nil {
		return l, nil
	}
	return l, writeError
}

// Lambda holds all the needed information for our lambda
type Lambda struct {
	Name    string
	Code    []byte
	Dir     string
	Custom  bool
	Command string
}

// Write takes the code and writes it to the directory + name
func (l *Lambda) Write() error {
	if err := ioutil.WriteFile(l.Dir+"/"+l.Name, l.Code, 0755); err == nil {
		return nil
	}
	return errors.New("Unable to write file")
}

// Execute will execute the lambda and return it's output and errors(if applicable)
func (l *Lambda) Execute(stdin io.Reader, args string) ([]byte, error) {
	args = strings.Replace(strings.TrimSpace(args), "/", " ", -1)
	if args != "" {
		args = " " + args
	}

	cmd := exec.Command("bash", "-c", l.Command+" "+l.Dir+"/"+l.Name+args)
	if l.Custom {
		cmd = exec.Command("bash", "-c", l.Command+args)
	}

	// pass through some stdin goodness
	cmd.Stdin = stdin

	// for those who are about to rock, I salute you.
	stdoutStderr, err := cmd.CombinedOutput()

	if err == nil {
		log.WithFields(log.Fields{"name": l.Name, "command": l.Command}).Info("Lambda Execution")
		return stdoutStderr, nil
	}

	log.WithFields(log.Fields{"name": l.Name, "command": l.Command}).Error("Lambda Execution")
	return stdoutStderr, errors.New("Error runing command")
}
