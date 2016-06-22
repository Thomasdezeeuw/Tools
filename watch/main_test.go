// Copyright (C) 2014 Thomas de Zeeuw.
// Licensed onder the MIT license that can be found in the LICENSE file.

// TODO: test throttling
package main

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

const instant = 0 * time.Microsecond

func TestMain(m *testing.M) {
	oldRunCmdTimeout := runCmdTimeout
	runCmdTimeout = 0 * time.Millisecond

	os.Exit(m.Run())

	runCmdTimeout = oldRunCmdTimeout
}

func TestNewCmd(t *testing.T) {
	t.Skip()
	expected := "Hi\n"
	cmds := []string{"echo", expected}
	timeout := instant

	cmd := Cmd{
		Cmd:     cmds,
		Quiet:   false,
		Timeout: timeout,
	}

	if len(cmd.Cmd) != len(cmds) {
		t.Errorf("Epexcted the command to be called to be %v, but got %v", cmds, cmd.Cmd)
	} else {
		for i, c := range cmd.Cmd {
			if c != cmds[i] {
				t.Errorf("Epexcted the command to be called to be %v, but got %v", cmds, cmd.Cmd)
				return // Running the command won't go correctly
			}
		}
	}

	if cmd.Timeout != timeout {
		t.Errorf("Expected the timeout to be %s, got %s", cmd.Timeout, timeout)
	}

	if cmd.running != false {
		t.Error("The command is running, but we didn't start it")
	}

	if cmd.Quiet != false {
		t.Error("Epexcted the command to output, but it's quiet")
		return // Running the command won't return anything, no point it testing it
	}

	got := getOutput(cmd, false)

	if got != expected {
		t.Errorf("Epexcted %s to be echoed, got %s", expected, got)
	}
}

// TODO: proper testing, now it tests if nothing comes up..
func TestEmptyCmd(t *testing.T) {
	cmd := Cmd{
		Cmd:     []string{},
		Quiet:   false,
		Timeout: instant,
	}

	expected := ""
	got := getOutput(cmd, false)

	if got != expected {
		t.Errorf("Epexcted %s to be echoed, got %s", expected, got)
	}
}

func TestErrorCmd(t *testing.T) {
	cmd := Cmd{
		Cmd:     []string{"SOME_ERROR_COMMAND"},
		Quiet:   false,
		Timeout: instant,
	}

	expected := "exec: \"SOME_ERROR_COMMAND\": executable file not found in %PATH%\n"
	got := getOutput(cmd, true)

	if got != expected {
		t.Errorf("Epexcted %s to be echoed, got %s", expected, got)
	}
}

func getOutput(cmd Cmd, errOutput bool) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	var oldStdout *os.File
	if !errOutput {
		oldStdout = os.Stdout
		os.Stdout = w
	} else {
		oldStdout = os.Stderr
		os.Stderr = w
	}

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)

		out <- buf.String()
	}()

	cmd.Run()

	w.Close()
	if !errOutput {
		os.Stdout = oldStdout
	} else {
		os.Stderr = oldStdout
	}

	return <-out
}
