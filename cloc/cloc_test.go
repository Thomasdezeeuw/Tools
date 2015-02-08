// Copyright (C) 2015 Thomas de Zeeuw.
//
// Licensed onder the MIT license that can be found in the LICENSE file.

// TODO(Thomas): Add code to files with 0 lines of code in the test directory
// TODO(Thomas): Test with errors.
package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestGetLanguage(t *testing.T) {
	type test struct {
		filepath string
		expected *language
	}

	tests := []test{
		{"file.as", languages["actionscript"]},
		{"file.asa", languages["asp"]},
		{"file.asp", languages["asp"]},
		{"file.c", languages["c"]},
		{"file.h", languages["c"]},
		{"file.cs", languages["c#"]},
		{"file.c++", languages["c++"]},
		{"file.cpp", languages["c++"]},
		{"file.cp", languages["c++"]},
		{"file.cc", languages["c++"]},
		{"file.hh", languages["c++"]},
		{"file.clj", languages["clojure"]},
		{"file.css", languages["css"]},
		{"file.d", languages["d"]},
		{"file.di", languages["d"]},
		{"file.erl", languages["erlang"]},
		{"file.hrl", languages["erlang"]},
		{"file.go", languages["go"]},
		{"file.dot", languages["dot"]},
		{"file.DOT", languages["dot"]},
		{"file.groovy", languages["groovy"]},
		{"file.gvy", languages["groovy"]},
		{"file.hs", languages["haskell"]},
		{"file.html", languages["html"]},
		{"file.htm", languages["html"]},
		{"file.shtml", languages["html"]},
		{"file.xhtml", languages["html"]},
		{"file.phtml", languages["html"]},
		{"file.tmpl", languages["html"]},
		{"file.tpl", languages["html"]},
		{"file.java", languages["java"]},
		{"file.js", languages["javascript"]},
		{"file.jsx", languages["javascript"]},
		{"file.lisp", languages["lisp"]},
		{"file.cl", languages["lisp"]},
		{"file.l", languages["lisp"]},
		{"file.lua", languages["lua"]},
		{"file.m", languages["objective-c"]},
		{"file.mm", languages["objective-c"]},
		{"file.M", languages["objective-c"]},
		{"file.ml", languages["ocaml"]},
		{"file.mli", languages["ocaml"]},
		{"file.mll", languages["ocaml"]},
		{"file.pas", languages["pascal"]},
		{"file.p", languages["pascal"]},
		{"file.pl", languages["perl"]},
		{"file.pm", languages["perl"]},
		{"file.php", languages["php"]},
		{"file.py", languages["python"]},
		{"file.rpy", languages["python"]},
		{"file.cpy", languages["python"]},
		{"file.pyw", languages["python"]},
		{"file.r", languages["r"]},
		{"file.R", languages["r"]},
		{"file.s", languages["r"]},
		{"file.S", languages["r"]},
		{"file.rb", languages["ruby"]},
		{"file.rbx", languages["ruby"]},
		{"file.rjs", languages["ruby"]},
		{"file.rs", languages["rust"]},
		{"file.scala", languages["scala"]},
		{"file.sh", languages["shell"]},
		{"file.bash", languages["shell"]},
		{"file.zsh", languages["shell"]},
		{"file.sql", languages["sql"]},
		{"file.txt", languages["unkown"]},
		{"somefile", languages["unkown"]},
		{"go", languages["unkown"]},
	}

	for _, test := range tests {
		lang := getLanguage(test.filepath)

		if lang != test.expected {
			t.Errorf("Expected getLanguage(%s) to return %v, got %v",
				test.filepath, test.expected, lang)
		}
	}
}

func TestGetFileOptions(t *testing.T) {
	type test struct {
		args     []string
		expected []string
	}

	tests := []test{
		{[]string{"bin"}, []string{"./"}},
		{[]string{"bin", "test"}, []string{"test"}},
		{[]string{"bin", "test", "test.go"}, []string{"test", "test.go"}},
		{[]string{"bin", "", "test", "test.go"}, []string{"test", "test.go"}},
		{[]string{"bin", "-d", "test", "test.go"}, []string{"test.go"}},
		{[]string{"bin", "--d", "test", "test.go"}, []string{"test.go"}},
		{[]string{"bin", "-", "test", "test.go"}, []string{"test.go"}},
		{[]string{"bin", "--", "test", "test.go"}, []string{"test.go"}},
	}

	for _, test := range tests {
		result := getFileOptions(test.args)

		if len(result) != len(test.expected) {
			t.Errorf("Expected getFileOptions(%v) to return %v, but got %v",
				test.args, test.expected, result)
			return
		}

		for i, file := range result {
			if test.expected[i] != file {
				t.Errorf("Expected getFileOptions(%v) to return %v, but got %v",
					test.args, test.expected, result)
			}
		}
	}
}

func TestCountFile(t *testing.T) {
	type test struct {
		filepath      string
		expectedCount int
		err           string
	}

	tests := []test{
		{"file", 0, ""},
		{"file.as", 0, ""},
		{"file.asp", 0, ""},
		{"file.c", 13, ""},
		{"file.cs", 0, ""},
		{"file.go", 9, ""},
		{"file.gvy", 0, ""},
		{"file.h", 6, ""},
		{"file.hs", 0, ""},
		{"file.htm", 10, ""},
		{"file.l", 0, ""},
		{"file.php", 1, ""},
		{"README.md", 0, ""},
		{"testdata1/file.cl", 0, ""},
		{"testdata1/file.clj", 0, ""},
		{"testdata1/file.cpp", 13, ""},
		{"testdata1/file.lisp", 0, ""},
		{"testdata1/file.lua", 0, ""},
		{"testdata1/file.m", 0, ""},
		{"testdata1/file.p", 0, ""},
		{"testdata1/file.txt", 0, ""},
		{"testdata1/README.md", 0, ""},
		{"testdata2/file.css", 3, ""},
		{"testdata2/file.d", 0, ""},
		{"testdata2/file.dot", 0, ""},
		{"testdata2/file.erl", 0, ""},
		{"testdata2/file.html", 10, ""},
		{"testdata2/file.java", 0, ""},
		{"testdata2/file.js", 5, ""},
		{"testdata2/README.md", 0, ""},
		{"testdata2/testdata3/file.pl", 1, ""},
		{"testdata2/testdata3/file.py", 1, ""},
		{"testdata2/testdata3/file.r", 0, ""},
		{"testdata2/testdata3/file.rb", 5, ""},
		{"testdata2/testdata3/file.rs", 15, ""},
		{"testdata2/testdata3/file.scala", 0, ""},
		{"testdata2/testdata3/file.sh", 1, ""},
		{"testdata2/testdata3/file.tmpl", 0, ""},
		{"testdata2/testdata3/README.md", 0, ""},
		{"not_found", 0, ""},
		{"not_found.go", 0, "open _testdata" + string(os.PathSeparator) +
			"not_found.go: The system cannot find the file specified."},
	}

	for _, test := range tests {
		count, err := countFile("_testdata/" + test.filepath)

		if err != nil && err.Error() != test.err {
			t.Errorf("Unexpected error %s, expected %s", err.Error(), test.err)
			continue
		}

		if count != test.expectedCount {
			t.Errorf("Expected %d lines of code, but got %d for file %s",
				test.expectedCount, count, test.filepath)
		}
	}
}

func TestCountDir(t *testing.T) {
	type test struct {
		filepath      string
		expectedCount int
		err           string
	}

	tests := []test{
		{"_testdata", 93, ""},
		{"not_found", 0, "open not_found: The system cannot find the file specified."},
	}

	for _, test := range tests {
		count, err := countDir(test.filepath)

		if err != nil && err.Error() != test.err {
			t.Errorf("Unexpected error %s, expected %s", err.Error(), test.err)
			return
		}

		if count != test.expectedCount {
			t.Errorf("Expected %d lines of code, but got %d for directory %s",
				test.expectedCount, count, test.filepath)
		}
	}
}

func TestCount(t *testing.T) {
	type test struct {
		filepath      string
		expectedCount int
		err           string
	}

	tests := []test{
		{"_testdata/file.go", 9, "nil"},
		{"_testdata", 93, "nil"},
		{"notFound", 0, "Cannot open file notFound."},
		{"notFound.go", 0, "Cannot open file notFound.go."},
	}

	for _, test := range tests {
		count, err := count(test.filepath)

		if err != nil && err.Error() != test.err {
			t.Errorf("Expected error to be \"%v\", got \"%v\"", test.err, err)
			return
		}

		if count != test.expectedCount {
			t.Errorf("Expected %d lines of code, but got %d for file/directory %s",
				test.expectedCount, count, test.filepath)
		}
	}
}

func TestMain(t *testing.T) {
	type test struct {
		args     []string
		expected string
	}

	tests := []test{
		{[]string{"", "languages.go"}, "languages.go contains 162 lines of code.\n"},
		{[]string{"", "languages.go", "README.md"}, "languages.go contains 162 lines of code.\n" +
			"README.md contains 0 lines of code.\nTotal number of lines: 162.\n"},
	}

	oldStdout := os.Stdout
	oldArgs := os.Args

	for _, test := range tests {
		os.Args = test.args

		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}

		os.Stdout = w
		outputChannel := make(chan string, len(tests))

		go func(i int) {
			for i > 0 {
				var buf bytes.Buffer
				io.Copy(&buf, r)

				outputChannel <- buf.String()
				i--
			}
		}(len(tests))

		main()

		w.Close()

		output := <-outputChannel

		if output != test.expected {
			t.Errorf("Expected the output to be '%s', got '%s'", test.expected, output)
		}
	}

	os.Stdout = oldStdout
	os.Args = oldArgs
}
