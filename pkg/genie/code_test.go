package genie

import (
	"io"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func example(stdin io.Reader, args []string) (string, error) {
	r, _ := ioutil.ReadAll(stdin)
	return "custom:" + string(r) + strings.Join(args, " "), nil
}

func TestCodeLambda(t *testing.T) {
	l := NewCodeLambda("example", example)

	out, _ := l.Execute(strings.NewReader("abcd"), []string{"efghi"})

	if out != "custom:abcdefghi" {
		log.Fatalf("Execute() on custom code should return stdin + args")
	}
}
