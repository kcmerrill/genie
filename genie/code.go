package genie

import "io"

type execute func(stdin io.Reader, args string) (string, error)

// NewCodeLambda will create a new code lambda
func NewCodeLambda(name string, fn execute) *CodeLambda {
	return &CodeLambda{name: name, fn: fn}
}

// CodeLambda will execute custom go code
type CodeLambda struct {
	name string
	fn   func(stdin io.Reader, args string) (string, error)
}

// Name returns the lambdas name
func (l *CodeLambda) Name() string {
	return l.name
}

// Execute executes the function and returns the results
func (l *CodeLambda) Execute(stdin io.Reader, args string) (string, error) {
	return l.fn(stdin, args)
}
