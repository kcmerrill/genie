package genie

import (
	"strings"
	"testing"
)

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

func TestGitHubLambda(t *testing.T) {
	g := New("/tmp", "port", "token")
	g.GithubLambda("github", "kcmerrill", "genie", "lambdas/echo.py")

	if g.Lambdas["github"].Dir != "/tmp" {
		t.Fatal("Github lambda directory is incorrect")
	}

	if g.Lambdas["github"].Command != "python" {
		t.Fatal("Github lambda command should be python")
	}

	if g.Lambdas["github"].Name != "github" {
		t.Fatal("Github lambda name should be github")
	}

	// silly hack, but a test is a test
	if !strings.HasPrefix(string(g.Lambdas["github"].Code), "import sys") {
		t.Fatal("Github lambda code should start with 'import sys'")
	}
}
