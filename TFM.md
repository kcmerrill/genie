# TFM

Welcome to *coughs* the manual.

Genie was designed to run simple repeatable bits of code. I like to call it WaaS(Whatever as a service) because it's really whatever you want to be a service to be. It could be a full fledge app, or a simple ls or a simple echo in whatever language you want.

You can install custom commands, you can use public github files, enter your own, again, it's really up to you as to what you can use it for.

## Usage

We'll assume `genie` is running at `http://localhost/`.

To register a new github public lambda named `my.new.python.lambda`:

```bash
curl -X GET http://localhost/my.new.python.lambda/github.com/kcmerrill/genie/lambdas/echo.py
```

To execute the lambda `my.new.python.lambda`

```bash
 curl -X GET http://localhost/my.new.python.lambda
```

To execute the lambda `my.new.python.lambda` with cli arguments

```bash
 curl -X GET http://localhost/my.new.python.lambda/cliarg1/cliarg2/etc/etc/etc
```

To execute the lambda `my.new.python.lambda` with stdin

```bash
curl -X POST -d "This goes straight to stdin\!" http://localhost:8080/my.new.python.lambda/arg1/arg2
```

To create a lambda via http body

```bash
curl -X POST -d "print 'My new lambda'" http://localhost:8080/mynewlambda/register/python
```

Where `mynewlambda` is the lambda name, `python` is the command to execute it.

## Custom code

You can use genie with custom go code as well. This is useful if you want to create lambdas with go code directly, along with all the other functionality you can get with genie and other programming languages.

Creating a new lambda is as simple as creating a new function with `func example(stdin io.Reader, args string) (string, error)`.
	
## Genie as a package

There were a few use cases I created `genie` for, not just for lambdas via the web browser. To use `genie` as a package, and install your own lambdas take a peek at `main.go`. In doing so, here are a few ways to inject your lambdas. 

```golang

NewCustomLambda(name, command string)
NewCodeLambda(name string, fn execute) // execute has the signature: func example(stdin io.Reader, args string) (string, error)
NewLocalLambda(name, directory, command string, code []byte) // where directory is genie.Dir(it has to know where to look)

```


## Options

When starting `genie` there are a few options you can use.

* `dir` will determine where to save your lambdas.
* `port` will determine which port to run on.
