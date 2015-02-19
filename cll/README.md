# Cll

[![Coverage](http://gocover.io/_badge/github.com/Thomasdezeeuw/tools/cll)](http://gocover.io/github.com/Thomasdezeeuw/tools/cll)

> Version 0.0.1

Cll checks if all lines have a length within the maximum allowed length.

## Installation

[Go](http://golang.org/) is required. Then run:

```bash
$ go get github.com/Thomasdezeeuw/tools/cll
$ cd $GOPATH/src/github.com/Thomasdezeeuw/tools/cll
$ go install
```

## Usage

Run through the command line:

```bash
$ cll
```

No result means no to long lines.

## Options

Specifying files and/or directories.

```bash
$ cloc my_file my_folder
my_folder/my_file line 2 is 90 characters long.
```

**Length (*-l*)**. The maxmimum allowed length of a line, defaults to 80.

```bash
$ cll -l 100
```

**Tab (*-t*)**. The length of a tab, defaults to 2

```bash
$ cll -t 4
```
