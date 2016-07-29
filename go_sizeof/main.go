package main

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/imports"
)

// TODO(Thomas): add a complete import path for the type, e.g.
// github.com/Thomasdezeeuw/ini.Config for ini.Conig. This is usefull making
// sure that the correct type is used.

func main() {
	flag.Parse()

	t := flag.Arg(0)
	if t == "" {
		flag.Usage()
		os.Exit(1)
	}

	path := filepath.Join(os.TempDir(), "tmp_main.go")

	if err := createSource(path, t); err != nil {
		exitErr("error creating source", err)
	}
	defer os.Remove(path)

	output, err := runSource(path)
	if err != nil {
		exitErr("error running go program", err)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		exitErr("error writing output to stdout", err)
	}
}

func createSource(path, t string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData{
		Type: t,
	})
	if err != nil {
		return err
	}

	source, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	_, err = f.Write(source)
	return err
}

func runSource(path string) ([]byte, error) {
	var errBuf, output bytes.Buffer

	cmd := exec.Command("go", "run", path)
	cmd.Stdout = &output
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		err = errors.New(err.Error() + "\n" + errBuf.String())
		return nil, err
	}

	return output.Bytes(), nil
}

func exitErr(extra string, err error) {
	os.Stderr.WriteString(extra + "\n")
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

var tmpl = template.Must(template.New("tmpl").Parse(templateText))

type tmplData struct {
	Type string
}

const templateText = `package main

	import (
		"fmt"
		"unsafe"
	)

	func main() {
		v := new({{.Type}})
		n := unsafe.Sizeof(*v)
		fmt.Printf("%d\n", n)
	}
`
