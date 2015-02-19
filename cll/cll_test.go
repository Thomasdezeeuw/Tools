// Copyright (C) 2015 Thomas de Zeeuw.
//
// Licensed onder the MIT license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

const (
	longLinesPath  = "_testdata" + string(os.PathSeparator) + "long-lines.txt"
	longLinesPath2 = "_testdata" + string(os.PathSeparator) + "testdata2" +
		string(os.PathSeparator) + "long-lines.txt"
)

func TestCheckFile(t *testing.T) {
	type test struct {
		input    string
		expected []longLine
		err      error
	}

	tests := []test{
		{"no-long-lines.txt", nil, nil},
		{"long-lines.txt", []longLine{{longLinesPath, 1, 566},
			{longLinesPath, 3, 551}, {longLinesPath, 5, 83}}, nil},
		/*{"not_found", nil, errors.New("open _testdata" + string(os.PathSeparator) +
		"not_found: The system cannot find the file specified.")},*/
	}

	for _, test := range tests {
		results, err := checkFile("_testdata/"+test.input, 80)

		errMsg := "Expected checkFile(_testdata/" + test.input + ", 80) "

		if err != test.err && fmt.Sprint("%v", err) != fmt.Sprint("%v", test.err) {
			t.Errorf(errMsg+"return error to be '%v', but got '%v'", test.err, err)
			continue
		} else if err != nil && test.err != nil {
			continue
		}

		if len(results) != len(test.expected) {
			t.Errorf(errMsg+"to return %v, but got %v", test.expected, results)
			continue
		}

		for i, result := range results {
			if result != test.expected[i] {
				t.Errorf(errMsg+"to return %v, but got %v", test.expected, results)
				break
			}
		}
	}
}

func TestCheckDir(t *testing.T) {
	type test struct {
		input    string
		expected []longLine
		err      error
	}

	tests := []test{
		{"testdata3", nil, nil},
		{"testdata2", []longLine{{longLinesPath2, 2, 205},
			{longLinesPath2, 3, 270}}, nil},
		{"", []longLine{{longLinesPath, 1, 566}, {longLinesPath, 3, 551},
			{longLinesPath, 5, 83}, {longLinesPath2, 2, 205},
			{longLinesPath2, 3, 270}}, nil},
		/*{"not_found", nil, errors.New("open _testdata" + string(os.PathSeparator) +
		"not_found: The system cannot find the file specified.")},*/
	}

	for _, test := range tests {
		results, err := checkDir("_testdata/"+test.input, 80)

		errMsg := "Expected checkDir(_testdata/" + test.input + ", 80) "

		if err != test.err && fmt.Sprint("%v", err) != fmt.Sprint("%v", test.err) {
			t.Errorf(errMsg+"return error to be '%v', but got '%v'", test.err, err)
			continue
		} else if err != nil && test.err != nil {
			continue
		}

		if len(results) != len(test.expected) {
			t.Errorf(errMsg+"to return %v, but got %v", test.expected, results)
			continue
		}

		for i, result := range results {
			if result != test.expected[i] {
				t.Errorf(errMsg+"to return %v, but got %v", test.expected, results)
				break
			}
		}
	}
}

func TestCheck(t *testing.T) {
	type test struct {
		input    string
		expected []longLine
		err      error
	}

	tests := []test{
		{"no-long-lines.txt", nil, nil},
		{"long-lines.txt", []longLine{{longLinesPath, 1, 566},
			{longLinesPath, 3, 551}, {longLinesPath, 5, 83}}, nil},
		{"testdata3", nil, nil},
		{"testdata2", []longLine{{longLinesPath2, 2, 205},
			{longLinesPath2, 3, 270}}, nil},
		{"", []longLine{{longLinesPath, 1, 566}, {longLinesPath, 3, 551},
			{longLinesPath, 5, 83}, {longLinesPath2, 2, 205},
			{longLinesPath2, 3, 270}}, nil},
		/*{"not_found", nil, errors.New("Cannot open file _testdata" +
		string(os.PathSeparator) + "not_found.")},*/
	}

	for _, test := range tests {
		results, err := check("_testdata/"+test.input, 80)

		errMsg := "Expected check(_testdata/" + test.input + ", 80) "

		if err != test.err && fmt.Sprint("%v", err) != fmt.Sprint("%v", test.err) {
			t.Errorf(errMsg+"return error to be '%v', but got '%v'", test.err, err)
			continue
		} else if err != nil && test.err != nil {
			continue
		}

		if len(results) != len(test.expected) {
			t.Errorf(errMsg+"to return %v, but got %v", test.expected, results)
			continue
		}

		for i, result := range results {
			if result != test.expected[i] {
				t.Errorf(errMsg+"to return %v, but got %v", test.expected, results)
				break
			}
		}
	}
}

func TestMain(t *testing.T) {
	type test struct {
		args     []string
		expected string
	}

	tests := []test{
		{[]string{"bin", "_testdata/testdata3"}, ""},
		{[]string{"bin", "_testdata/testdata2"}, longLinesPath2 + " line 2 is " +
			"205 characters long.\n" + longLinesPath2 + " line 3 is 270 characters" +
			" long.\n"},
		{[]string{"bin", "_testdata"}, longLinesPath + " line 1 is 566 characters" +
			" long.\n" + longLinesPath + " line 3 is 551 characters long.\n" +
			longLinesPath + " line 5 is 83 characters long.\n" + longLinesPath2 +
			" line 2 is 205 characters long.\n" + longLinesPath2 + " line 3 is 270 " +
			"characters long.\n"},
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
			t.Errorf("Expected the output to be '%s', got '%s'",
				test.expected, output)
		}
	}

	os.Stdout = oldStdout
	os.Args = oldArgs
}
