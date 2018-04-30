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
	l := NewCustomLambda("woot", "echo woot!")
	g.AddLambda(l)

	if _, exists := g.Lambdas["woot"]; !exists {
		t.Fatal("AddLambda() should add a lambda!")
	}
}

func TestGitHubLambda(t *testing.T) {
	g := New("/tmp", "port", "token")
	g.GithubLambda("github", "kcmerrill", "genie", "lambdas/echo.py")

}

func TestExecuteLambda(t *testing.T) {
	g := New("/tmp", "port", "token")
	l := NewCustomLambda("echo", "echo")
	g.AddLambda(l)

	out, err := g.Execute("echo", strings.NewReader(""), []string{"kcwashere-indeed"})
	if out != "kcwashere-indeed" {
		t.Fatalf("Expecting: 'kcwashere-indeed', Actual: '%s'", out)
	}

	if err != nil {
		t.Fatalf("No errors on Execute()")
	}
}
