package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestGetGitLog(t *testing.T) {
	t.Parallel()

	r, err := getGitLog("_testdata/repo")
	if err != nil {
		t.Fatalf("Unexpected error getting git log: %s", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Unexpected error reading git log output: %s", err)
	}

	got := string(b)
	expected := `hash: e0daf77
author: Thomas de Zeeuw
date: 2016-06-22T14:43:14+02:00
ref: HEAD -> master, tag: v0.1
title: init()
message: A simple test commit, safe to ignore.

==============================

`

	if got != expected {
		diff := pretty.Compare(expected, got)
		t.Fatalf("Got value different from the expected value: \n%s", diff)
	}
}

func TestParseGitLog(t *testing.T) {
	t.Parallel()

	var r strings.Reader
	for i, test := range commitInput {
		r.Reset(test)

		got, err := parseGitLog(&r)
		if err != nil {
			t.Fatalf("Unexpected error parsing testdata: %s", err)
		}

		if len(got) != 1 {
			diff := pretty.Compare([]Commit{expectedCommits[i]}, got)
			t.Fatalf("Got value differences from the expected value: \n%s", diff)
		}

		g, e := got[0], expectedCommits[i]
		if !reflect.DeepEqual(e, g) {
			diff := pretty.Compare(e, g)
			t.Fatalf("Got value different from the expected value: \n%s", diff)
		}
	}
}

func TestWriteChangelog(t *testing.T) {
	t.Parallel()

	const prefix = "# Changelog\n\n"

	var buf bytes.Buffer
	for i, test := range expectedCommits {
		buf.Reset()

		if err := writeChangelog([]Commit{test}, &buf); err != nil {
			t.Fatalf("Unexpected error writing changelog: %s", err)
		}

		got := buf.String()

		if !strings.HasPrefix(got, prefix) {
			t.Fatalf("Expected the changelog to have the prefix %s, but got %s",
				prefix, got)
		}

		got = got[len(prefix):]
		expected := expectedChangelog[i]

		headerPrefix := "## Master\n\n"
		if test.Tag != "" {
			headerPrefix = fmt.Sprintf("## %s\n\n", test.Tag)
		}

		expected = headerPrefix + expected

		if got != expected {
			diff := pretty.Compare(expected, got)
			t.Fatalf("Got value different from the expected value: \n%s", diff)
		}
	}
}

func TestWriteChangelogWithTag(t *testing.T) {
	t.Parallel()

	commits := []Commit{
		{
			Hash:   "308ae12",
			Author: "Thomas de Zeeuw",
			Date:   mustParseTime("2016-06-22T15:43:14+02:00"),
			Title:  "Version 0.1",
			Tag:    "v0.1",
		},
		{
			Hash:    "3a08ae3",
			Author:  "Thomas de Zeeuw",
			Date:    mustParseTime("2016-06-22T14:43:14+02:00"),
			Title:   "init()",
			Message: "A simple test commit, safe to ignore.",
		},
	}

	const expected = "# Changelog\n\n## v0.1\n\n" +
		" - **Version 0.1** (#308ae12) by *Thomas de Zeeuw*, on *22 Jun 2016 13:43:14 UTC*.\n" +
		" - **init()** (#3a08ae3) by *Thomas de Zeeuw*, on *22 Jun 2016 12:43:14 UTC*. <br/>\n*A simple test commit, safe to ignore*.\n"

	var buf bytes.Buffer

	if err := writeChangelog(commits, &buf); err != nil {
		t.Fatalf("Unexpected error writing changelog: %s", err)
	}

	got := buf.String()

	if got != expected {
		diff := pretty.Compare(expected, got)
		t.Fatalf("Got value different from the expected value: \n%s", diff)
	}
}
