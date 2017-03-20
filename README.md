# genie

[![Build Status](https://travis-ci.org/kcmerrill/genie.svg?branch=master)](https://travis-ci.org/kcmerrill/genie) [![Join the chat at https://gitter.im/kcmerrill/genie](https://badges.gitter.im/kcmerrill/genie.svg)](https://gitter.im/kcmerrill/genie?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![genie](https://raw.githubusercontent.com/kcmerrill/genie/master/assets/genie.jpg "genie")

Lambda knockoff. A [crush](http://github.com/kcmerrill/crush) companion.

## Binaries || Installation

[![MacOSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX")](http://go-dist.kcmerrill.com/kcmerrill/genie/mac/amd64) [![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux")](http://go-dist.kcmerrill.com/kcmerrill/genie/linux/amd64)

via go:

`$ go get -u github.com/kcmerrill/genie`

## Usage

Lambdas can only(right now) be fetched from public github repos. Feel free to use the lambdas stored in this public repository as examples. For these demos, we'll assume `genie` is running at `http://localhost/`.

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

## Options

When starting `genie` there are a few options you can use.

* `dir` will determine where to save your lambdas.
* `port` will determine which port to run on.