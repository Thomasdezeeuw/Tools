// Copyright (C) 2014 Thomas de Zeeuw.
// Licensed onder the MIT license that can be found in the LICENSE file.

// This tools watches for file changes and then executes a given command.
// TODO: ONLY HAVE TIME WHEN SPECIFIED, otherwise just let the proces live and
// kill it when a change has occurered
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"gopkg.in/fsnotify.v1"
)

type Cmd struct {
	// The command with arguments we're going to run.
	Cmd []string

	// Wether the command should output to stdout & stderr.
	Quiet bool

	// Timeout for the command.
	Timeout time.Duration

	// Wether or not a command is being run.
	running bool
}

// Run the command with a throttle for the next run.
func (c *Cmd) Run() {
	// If the command is being run we don't need to do it again.
	// Or if no command is going to be run.
	if c.running || len(c.Cmd) < 1 {
		return
	}

	// The command we're going to run.
	cmd := exec.Command(c.Cmd[0], c.Cmd[1:]...)
	cmd.Env = os.Environ()

	// Output stderr and possible stdout
	cmd.Stderr = os.Stderr
	if !c.Quiet {
		cmd.Stdout = os.Stdout
	}

	// Show that we're running a command, used as a throttling mechanism.
	c.running = true

	// Wait for the latest file changes.
	time.Sleep(runCmdTimeout)

	// Run the command and print possible errors.
	if err := cmd.Start(); err != nil {
		fmt.Print(err.Error() + "\n")
	}

	// Have a timeout to make sure the process doesn't run forever
	timeout := time.AfterFunc(c.Timeout, func() {
		fmt.Print("Long running command, killing it\n")
		// First send the kill signal.
		cmd.Process.Kill()

		// Then make sure it actually got killed.
		time.AfterFunc(500*time.Millisecond, func() {
			// And if it didn't.
			if !cmd.ProcessState.Exited() {
				// Force it.
				cmd.Process.Kill()
			}
		})
	})

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Print(err.Error() + "\n")
	}

	timeout.Stop()

	// We're ready for another command.
	c.running = false
}

var (
	// Wether the standard output should be ignored.
	quiet = false

	// The timeout for long running processes.
	// TODO(Thomas): Change to 10 seconds.
	timeout = 10 * time.Second

	// Timeout before actually running the command, this is usefull beause files
	// might change in rapid succession.
	runCmdTimeout = 100 * time.Millisecond
)

func init() {
	flag.BoolVar(&quiet, "q", false, "Only show the error output, not the regular output.")
	flag.DurationVar(&timeout, "t", timeout, "Time to wait after file change to run the command.")
}

func main() {
	flag.Parse()

	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	// Watch the working directory
	// TODO(Thomas): add option to specify the directory to watch
	// TODO(Thomas): add ignore option
	// TODO(Thomas): Add recursive watching;
	// - Search current dir for all folders
	// - Add each folders
	// - For each add event check if it's a folder if so add it to watch list
	err = watcher.Add("./")
	if err != nil {
		panic(err)
	}

	// Create the command from the command
	cmd := Cmd{
		Cmd:     flag.Args(),
		Quiet:   quiet,
		Timeout: timeout,
	}

	// Run the command, usefull we using it in combination with testing
	cmd.Run()

	// Now we're going to wait for file changes and possible errors.
	for {
		select {
		case <-watcher.Events:
			go cmd.Run()

		case err := <-watcher.Errors:
			fmt.Print(err.Error() + "\n")
		}
	}
}
