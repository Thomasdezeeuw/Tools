// Copyright (C) 2015 Thomas de Zeeuw.
//
// Licensed onder the MIT license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const (
	portDesc = "The port to listen on, defaults to 8000"
	dirDesc  = "The directory to serve, defaults to the working directory"
)

var (
	port string
	dir  string
)

func init() {
	flag.StringVar(&port, "port", "8000", portDesc)
	flag.StringVar(&port, "p", "8000", portDesc)
	flag.StringVar(&dir, "directory", "", dirDesc)
	flag.StringVar(&dir, "d", "", dirDesc)
}

func main() {
	flag.Parse()

	dir = filepath.Join("./", dir)

	nameDir := dir
	if nameDir == "." {
		nameDir = "current directory"
	}

	fmt.Printf("Serving directory %s, on port %s.\n", nameDir, port)
	err := http.ListenAndServe(":"+port, http.FileServer(http.Dir(dir)))
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
