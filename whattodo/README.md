# Whattodo

## Wondering what work needs to be done?

Look no further, with this tool you can easily find what needs to be done
within your [Go](http://golang.org/) project!

This tool looks through your [Go](http://golang.org/) source code and gets all
the todo items out of it. The supported formats are `todo:` and `todo(name):`.

## Usage
````
$ whattodo
  - [ ] Add this great feature (main.go, line 13).
```

```
$ whattodo my_dir
  - [ ] Add this great feature (my_dir/main.go, line 08).
```

## Installation

[Go](http://golang.org/) is required, then run:

```bash
$ go get github.com/Thomasdezeeuw/whattodo
$ cd $GOPATH/src/github.com/Thomasdezeeuw/whattodo
$ go install
```
