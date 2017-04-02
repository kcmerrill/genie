package genie

import (
	"os"
	"strings"
	"testing"
)

func TestCustomLambdaOk(t *testing.T) {
	nl := NewCustomLambda("name", "echo")
	out, _ := nl.Execute(os.Stdin, "kcwazhere")
	if strings.TrimSpace(out) != "kcwazhere" {
		t.Fatal("With a simple echo, was expecting kcwazhere to be displayed")
	}
}

func TestCustomLambdaFail(t *testing.T) {
	nl := NewCustomLambda("custom.failure.stuff", "asdf")
	_, err := nl.Execute(os.Stdin, "kcwazhere")
	if err == nil {
		t.Fatal("Was expecting an error given the command asdf does not exist")
	}
}
