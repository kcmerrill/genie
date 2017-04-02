package genie

import (
	"io"
)

// Lambda only requirement is to execute a command and return the output and error(if one)
type Lambda interface {
	Name() string
	Execute(stdin io.Reader, args string) (string, error)
}
