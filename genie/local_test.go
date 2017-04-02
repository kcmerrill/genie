package genie

import (
	"io/ioutil"
	"testing"
)

func TestNewLocalLambda(t *testing.T) {
	NewLocalLambda("bleh", "/tmp/", "python", []byte("hello world"))
	contents, _ := ioutil.ReadFile("/tmp/bleh")
	if string(contents) != "hello world" {
		t.Fatal("We should've written a file /tmp/bleh with 'hello world' inside")
	}
}
