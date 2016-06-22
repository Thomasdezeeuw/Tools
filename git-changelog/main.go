package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	outputDesc = "file to write changelog to, default to stdout"
	shorthand  = " (shorthand)"
)

var output string

func init() {
	flag.StringVar(&output, "o", output, outputDesc)
	flag.StringVar(&output, "output", output, outputDesc+shorthand)
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	// To writer to write the changelog to, default to stdout but using -o and
	// --output it can be changed to a file.
	var out io.Writer = os.Stdout
	if output != "" {
		f, err := os.Open(output)
		if err != nil {
			log.Fatal("Failed to open output file:", err.Error())
		}
		defer f.Close()

		out = f
	}

	logR, err := getGitLog("")
	if err != nil {
		log.Fatal("Failed to get git log:", err.Error())
	}

	commits, err := parseGitLog(logR)
	if err != nil {
		log.Fatal("Failed to parse git log:", err.Error())
	}

	if err := writeChangelog(commits, out); err != nil {
		log.Fatal("Failed to write changelog:", err.Error())
	}
}

const (
	logSeparator = "=============================="

	logFormat = `hash: %h
author: %an
date: %cI
ref: %D
title: %s
message: %b
` + logSeparator + `%n`
)

// getGitLog runs `git log` with the above defined logFormat.
func getGitLog(repoPath string) (io.Reader, error) {
	cmd := exec.Command("git", "log", "--format="+logFormat)

	// Change the working directory if not the current directory (default), mostly
	// used in testing.
	if repoPath != "" {
		cmd.Dir = repoPath
	}

	// Catch both the stdout and stderr. Cmd.Run() only returns the exit status,
	// which is not not enough to understand what is going on.
	var buf, errBuf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return nil, errors.New(errBuf.String() + err.Error())
	}

	return &buf, nil
}

// TODO(Thomas): make commits link to the commit in GitHub, behind a flag.
//
// URL is https://github.com/$user/$repo/commit/$hash

// writeChangelog writes a changelog to the provided writer, using the provided
// commits. It will create a new header for each tag, starting with `master` if
// the first commit doesn't have a tag.
func writeChangelog(commits []Commit, w io.Writer) error {
	// TODO(Thomas): maybe parse tags in the git log aswell and use them to create
	// headers based on the tags, for example ## v0.0.1.

	io.WriteString(w, "# Changelog\n\n")

	// No commits, no changelog.
	if len(commits) == 0 {
		return nil
	}

	// Start at master, if the latest commit doesn't have a tag.
	if commits[0].Tag == "" {
		io.WriteString(w, "## Master\n\n")
	}

	for _, commit := range commits {
		if commit.Tag != "" {
			fmt.Fprintf(w, "## %s\n\n", commit.Tag)
		}

		title := strings.TrimSuffix(commit.Title, ".")

		// TODO(Thomas): different time format?
		fmt.Fprintf(w, " - **%s** (#%s) by *%s*, on *%s*.",
			title, commit.Hash, commit.Author, commit.Date.UTC().Format("02 Jan 2006 15:04:05 MST"))

		if commit.Message != "" {
			msg := strings.TrimSuffix(commit.Message, ".")

			fmt.Fprintf(w, " <br/>\n*%s*.\n", msg)
		} else {
			io.WriteString(w, "\n")
		}
	}

	return nil
}
