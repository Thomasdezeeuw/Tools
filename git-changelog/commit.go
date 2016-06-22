package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"time"
)

type Commit struct {
	Hash    string
	Author  string
	Date    time.Time
	Title   string
	Message string
	Tag     string // Optional tag, e.g. v0.1.
}

// Parses a git log with the following format:
//
//	hash: $commit_hash
//	author: $authour_name
//	date: $commit_date
//	title: $commit_title
//	message: $commit_message
//
// For example:
//
//	hash: b6652b0
//	author: Thomas de Zeeuw
//	date: 2015-04-20T00:12:13+02:00
//	title: Serve: Allow first argument as directory
//	message: Now you can call `serve dir` or `server -d dir`.
//
// When it reaches a parsing error it will return the error and the results
// parsed until the error was encounted.
func parseGitLog(r io.Reader) ([]Commit, error) {
	var commits []Commit

	s := bufio.NewScanner(r)
	s.Split(splitSingleCommit)

	for s.Scan() {
		b := s.Bytes()
		commit, err := parseLogCommit(b)
		if err != nil {
			return commits, err
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// Fields used in parseLogCommit.
//
// IF MODIFIED ALSO MODIFIY parseLogCommit SWITCH!
var fields = [...][]byte{[]byte("hash"), []byte("author"), []byte("date"), []byte("ref"), []byte("title")}

// Parses a single git log commit, with the following format:
//
//	hash: $commit_hash
//	author: $authour_name
//	date: $commit_date
//	title: $commit_title
//	message: $commit_message
//
// For example:
//
//	hash: b6652b0
//	author: Thomas de Zeeuw
//	date: 2015-04-20T00:12:13+02:00
//	title: Serve: Allow first argument as directory
//	message: Now you can call `serve dir` or `server -d dir`.
func parseLogCommit(b []byte) (Commit, error) {
	var commit Commit

	for i, field := range fields {
		advance, value, err := getValue(b, field)
		if err != nil {
			return Commit{}, err
		}

		switch i {
		case 0:
			commit.Hash = string(value)
		case 1:
			commit.Author = string(value)
		case 2:
			t, err := time.Parse(time.RFC3339, string(value))
			if err != nil {
				return Commit{}, err
			}
			commit.Date = t.UTC()
		case 3:
			tag, ok := parseTag(value)
			if ok {
				commit.Tag = tag
			}
		case 4:
			commit.Title = string(value)
		}

		b = b[advance:]
	}

	message, err := getCommitMessage(b)
	if err != nil {
		return Commit{}, err
	}

	commit.Message = string(message)

	return commit, nil
}

const (
	separator byte = ':'
	newLine   byte = '\n'
	space     byte = ' '
)

// getValue splits a line and returns the value, using the following format:
//
//	field: value
//
// For example:
//
//	hash: b6652b0
func getValue(b, field []byte) (advance int, value []byte, err error) {
	i := bytes.IndexByte(b, separator)
	if i <= 0 {
		// TODO(Thomas): improve error message, this means nothing to the end user.
		return 0, []byte{}, errors.New("No separator found")
	}

	// TODO(Thomas): test field.

	// Drop: "field:"
	line := b[i+1:]
	advance = i + 1

	// Value ends at the end of the line, so advance the bytes atleast until then.
	i = bytes.IndexByte(line, newLine)
	advance = advance + i
	if i > 0 {
		// Slice until the end of the line.
		line = line[:i]
	}

	value = bytes.TrimSpace(line)
	return advance, value, nil
}

// parseTag
//
//	ref: HEAD -> master, origin/master, tag: v0.2
func parseTag(b []byte) (tag string, ok bool) {
	var tagPrefix = []byte("tag: ")

	//  Ref contains a bunch of information, we only need the tag part; `tag:
	//  $tag`. So we'll check if it's present and if so slice every before off the
	//  b slice.
	i := bytes.Index(b, tagPrefix)
	if i == -1 {
		return "", false
	}
	b = b[i+len(tagPrefix):]

	// Now we'll check if there's more information after the tag will also remove.
	if i := bytes.IndexByte(b, ','); i > 0 {
		b = b[:i]
	}

	tag = string(bytes.TrimSpace(b))
	return tag, true
}

// getCommitMessage parses a commit message, using the following format:
//
//	message: value
//
// For example:
//
//	message: Now you can call `serve dir` or `server -d dir`.
//
// The message may contain multiple lines and extra new lines at the end of it.
func getCommitMessage(b []byte) (value []byte, err error) {
	i := bytes.IndexByte(b, separator)
	if i <= 0 {
		// TODO(Thomas): improve error message, this means nothing to the end user.
		return []byte{}, errors.New("No separator found")
	}

	// TODO(Thomas): test field.

	// Drop: "message:"
	line := b[i+1:]

	// The commit message might be mulitple lines. To make a single line of this
	// we'll split on newlines, then trim any excess space of the lines and join
	// them using spaces.
	allLines := bytes.Split(line, []byte{newLine})
	for i, line := range allLines {
		allLines[i] = bytes.TrimSpace(line)
	}
	value = bytes.Join(allLines, []byte{space})
	value = bytes.TrimSpace(value)

	return value, nil
}

var commitSplit = []byte(logSeparator)

func splitSingleCommit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, commitSplit); i >= 0 {
		return i + len(commitSplit), data[:i], nil
	}

	// Request more data.
	return 0, nil, nil
}
