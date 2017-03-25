package genie

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewLambda(t *testing.T) {
	nl, _ := NewLambda("/tmp", "bleh", "python", []byte("hello world"))

	if nl.Dir != "/tmp" {
		t.Fatal("nl.dir should = /tmp")
	}

	if nl.Name != "bleh" {
		t.Fatal("The lambda's name should be bleh")
	}

	if string(nl.Code) != "hello world" {
		t.Fatal("The lambdas code should = 'hello world'")
	}

	contents, _ := ioutil.ReadFile("/tmp/bleh")
	if string(contents) != "hello world" {
		t.Fatal("We should've written a file /tmp/bleh with 'hello world' inside")
	}
}

func TestCustomLambdaOk(t *testing.T) {
	nl, _ := NewLambda("/tmp", "custom.stuff", "echo", []byte(""))
	nl.Custom = true
	out, _ := nl.Execute(os.Stdin, "kcwazhere")
	if strings.TrimSpace(string(out)) != "kcwazhere" {
		t.Fatal("With a simple echo, was expecting kcwazhere to be displayed")
	}
}

func TestCustomLambdaFail(t *testing.T) {
	nl, _ := NewLambda("/tmp", "custom.failure.stuff", "asdf", []byte(""))
	nl.Custom = true
	_, err := nl.Execute(os.Stdin, "kcwazhere")
	if err == nil {
		t.Fatal("Was expecting an error given the command asdf does not exist")
	}
}
