package genie

import (
	"errors"
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
func (l *Lambda) Execute(args string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", l.Command+" "+l.Dir+"/"+l.Name+" "+strings.Replace(args, "/", " ", -1))
	stdoutStderr, err := cmd.CombinedOutput()
	if err == nil {
		log.WithFields(log.Fields{"name": l.Name, "command": l.Command}).Info("Lambda Ran(succesful)")
		return stdoutStderr, nil
	}
	log.WithFields(log.Fields{"name": l.Name, "command": l.Command}).Error("Lambda Ran(failed)")
	return stdoutStderr, errors.New("Error runing command")
}
