# Cloc

[![Coverage](http://gocover.io/_badge/github.com/Thomasdezeeuw/tools/cloc)](http://gocover.io/github.com/Thomasdezeeuw/tools/cloc)

> Version 0.0.1

Cloc counts the number of lines of code in a given file or directory.

## Installation

[Go](http://golang.org/) is required. Then run:

```bash
$ go get github.com/Thomasdezeeuw/tools/cloc
$ cd $GOPATH/src/github.com/Thomasdezeeuw/tools/cloc
$ go install
```

## Usage

Run through the command line:

```bash
$ cloc
Current directory contains 291 lines of code.
```

## Options

Specifying files and/or directories.

```bash
$ cloc my_file my_folder
my_file contains 60 lines of code.
my_folder contains 100 lines of code.
Total number of lines: 160.
```
