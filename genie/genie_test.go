package genie

import "testing"

func TestNew(t *testing.T) {
	g := New("dir", "port", "token")
	if g.Dir != "dir" {
		t.Fatal("Expecting g.dir = 'dir'")
	}

	if g.Port != "port" {
		t.Fatal("Execting g.port = 'port'")
	}

	if g.Token != "token" {
		t.Fatal("Expecting g.token = 'token'")
	}

	if len(g.Lambdas) != 0 {
		t.Fatal("g.Lambas should be empty")
	}
}

func TestAddLambda(t *testing.T) {
	g := New("dir", "port", "token")
	l := &Lambda{
		Name:    "woot",
		Dir:     "dir",
		Code:    []byte("hello world"),
		Command: "command",
	}

	g.AddLambda(l)

	if _, exists := g.Lambdas["woot"]; !exists {
		t.Fatal("AddLambda() should add a lambda!")
	}
}

func TestGenerateCommand(t *testing.T) {
	g := New("dir", "port", "token")
	// test php
	if g.GenerateCommand("file.php") != "php" {
		t.Fatal(".php should return php")
	}

	// test python
	if g.GenerateCommand("file.py") != "python" {
		t.Fatal(".py should return python")
	}

	// test ruby
	if g.GenerateCommand("file.rb") != "ruby" {
		t.Fatal(".rb should return ruby")
	}

	// test sh
	if g.GenerateCommand("file.sh") != "sh" {
		t.Fatal(".sh should return sh")
	}

	// test default
	if g.GenerateCommand("file.doesnotexist") != "" {
		t.Fatal("Unknown filetypes should return an empty string")
	}
}
