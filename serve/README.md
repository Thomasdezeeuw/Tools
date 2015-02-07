# Serve

> Version 0.0.1

Simply serve static files from a given directory on a given port.

## Installation

[Go](http://golang.org/) is required. Then run:

```bash
$ go get github.com/Thomasdezeeuw/tools/serve
$ cd $GOPATH/src/github.com/Thomasdezeeuw/tools/serve
$ go install
```

## Usage

Run through the command line:

```bash
$ serve
```

## Options

**Directory (*-d*)**. What directory to serve, defaults to the working directory.

```bash
$ serve -d public
```

**Port (*-p*)**. The port to listen on, defaults to 8000.

```bash
$ serve -p 9000
```
