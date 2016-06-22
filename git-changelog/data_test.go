package main

import "time"

func mustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}

// Must match up with expectedChangelog and commitInput below.
var expectedCommits = []Commit{
	{
		Hash:    "b6652b0",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-04-20T00:12:13+02:00"),
		Title:   "Serve: Allow first argument as directory",
		Message: "Now you can call `serve dir` or `server -d dir`.",
		Tag:     "v0.2",
	},
	{
		Hash:    "7c3924c",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T20:33:13+01:00"),
		Title:   "Move Gocover badges to the subdirectories",
		Message: "",
	},
	{
		Hash:    "f7705fc",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T20:31:43+01:00"),
		Title:   "Drop godoc badge from main readme",
		Message: "",
	},
	{
		Hash:    "5f5e82f",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T20:27:07+01:00"),
		Title:   "all: adhere to 80 character line limit",
		Message: "",
	},
	{
		Hash:    "14600c3",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T20:25:33+01:00"),
		Title:   "cloc: skip file not found tests",
		Message: "same as previous commit, but these changes didn't get commited.",
	},
	{
		Hash:    "14d54f7",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T20:24:08+01:00"),
		Title:   "cll, cloc: skip file not found tests",
		Message: "The tests fail on travi-ci.org. The problem is likely that the error message comes from the os and are therefor not platform independent.",
	},
	{
		Hash:    "143d198",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T20:02:08+01:00"),
		Title:   "Fix Gocover badge",
		Message: "",
		Tag:     "v0.1",
	},
	{
		Hash:    "0c05371",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T19:59:10+01:00"),
		Title:   "Add Travis, Gocover and Godoc badges",
		Message: "",
	},
	{
		Hash:    "fc75e22",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T19:56:38+01:00"),
		Title:   "Fix travis config",
		Message: ".travis.yml requires spaces, not tabs.",
	},
	{
		Hash:    "99d43ae",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T19:52:55+01:00"),
		Title:   "Add editorconfig",
		Message: "",
	},
	{
		Hash:    "843f6e4",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-19T19:52:46+01:00"),
		Title:   "Add Travis-ci config",
		Message: "Adds testing for go 1.3, 1.4 and tip.",
	},
	{
		Hash:    "8188c6f",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-08T15:59:54+01:00"),
		Title:   "All: update examples in readme",
		Message: "",
	},
	{
		Hash:    "141942b",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-08T02:53:56+01:00"),
		Title:   "Cloc: improve readme example",
		Message: "`cloc my_file my_folder` is clearer then `cloc cloc.go _testdata`",
	},
	{
		Hash:    "1bfdcf0",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-08T02:53:05+01:00"),
		Title:   "Add Cll",
		Message: "Cll checks if all lines have a length within the maximum allowed length.",
	},
	{
		Hash:    "545aa9d",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-08T01:29:41+01:00"),
		Title:   "All: update copyright notice in each source file",
		Message: "",
	},
	{
		Hash:    "c6305e9",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-08T00:37:11+01:00"),
		Title:   "Serve: write errors to stderr instead of stdout",
		Message: "",
	},
	{
		Hash:    "f5e25e4",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-08T00:36:41+01:00"),
		Title:   "Cloc: write errors to stderr instead of stdout",
		Message: "",
	},
	{
		Hash:    "d557a6c",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-07T15:23:13+01:00"),
		Title:   "Cloc: clean path in every count function",
		Message: "",
	},
	{
		Hash:    "c4ab7c7",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-07T15:05:31+01:00"),
		Title:   "Cloc: stop after we encounter an error",
		Message: "",
	},
	{
		Hash:    "9014363",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-07T02:44:54+01:00"),
		Title:   "Add cloc",
		Message: "Cloc counts the number of lines of code in a given file or directory",
	},
	{
		Hash:    "a7a6823",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-07T02:44:33+01:00"),
		Title:   "Add serve",
		Message: "Simply serve static files from a given directory on a given port",
	},
	{
		Hash:    "7539f40",
		Author:  "Thomas de Zeeuw",
		Date:    mustParseTime("2015-02-07T02:43:58+01:00"),
		Title:   "init()",
		Message: "",
	},
}

// Must match up with expectedCommits above and commitInput below.
var expectedChangelog = []string{
	" - **Serve: Allow first argument as directory** (#b6652b0) by *Thomas de Zeeuw*, on *19 Apr 2015 22:12:13 UTC*. <br/>\n*Now you can call `serve dir` or `server -d dir`*.\n",
	" - **Move Gocover badges to the subdirectories** (#7c3924c) by *Thomas de Zeeuw*, on *19 Feb 2015 19:33:13 UTC*.\n",
	" - **Drop godoc badge from main readme** (#f7705fc) by *Thomas de Zeeuw*, on *19 Feb 2015 19:31:43 UTC*.\n",
	" - **all: adhere to 80 character line limit** (#5f5e82f) by *Thomas de Zeeuw*, on *19 Feb 2015 19:27:07 UTC*.\n",
	" - **cloc: skip file not found tests** (#14600c3) by *Thomas de Zeeuw*, on *19 Feb 2015 19:25:33 UTC*. <br/>\n*same as previous commit, but these changes didn't get commited*.\n",
	" - **cll, cloc: skip file not found tests** (#14d54f7) by *Thomas de Zeeuw*, on *19 Feb 2015 19:24:08 UTC*. <br/>\n*The tests fail on travi-ci.org. The problem is likely that the error message comes from the os and are therefor not platform independent*.\n",
	" - **Fix Gocover badge** (#143d198) by *Thomas de Zeeuw*, on *19 Feb 2015 19:02:08 UTC*.\n",
	" - **Add Travis, Gocover and Godoc badges** (#0c05371) by *Thomas de Zeeuw*, on *19 Feb 2015 18:59:10 UTC*.\n",
	" - **Fix travis config** (#fc75e22) by *Thomas de Zeeuw*, on *19 Feb 2015 18:56:38 UTC*. <br/>\n*.travis.yml requires spaces, not tabs*.\n",
	" - **Add editorconfig** (#99d43ae) by *Thomas de Zeeuw*, on *19 Feb 2015 18:52:55 UTC*.\n",
	" - **Add Travis-ci config** (#843f6e4) by *Thomas de Zeeuw*, on *19 Feb 2015 18:52:46 UTC*. <br/>\n*Adds testing for go 1.3, 1.4 and tip*.\n",
	" - **All: update examples in readme** (#8188c6f) by *Thomas de Zeeuw*, on *08 Feb 2015 14:59:54 UTC*.\n",
	" - **Cloc: improve readme example** (#141942b) by *Thomas de Zeeuw*, on *08 Feb 2015 01:53:56 UTC*. <br/>\n*`cloc my_file my_folder` is clearer then `cloc cloc.go _testdata`*.\n",
	" - **Add Cll** (#1bfdcf0) by *Thomas de Zeeuw*, on *08 Feb 2015 01:53:05 UTC*. <br/>\n*Cll checks if all lines have a length within the maximum allowed length*.\n",
	" - **All: update copyright notice in each source file** (#545aa9d) by *Thomas de Zeeuw*, on *08 Feb 2015 00:29:41 UTC*.\n",
	" - **Serve: write errors to stderr instead of stdout** (#c6305e9) by *Thomas de Zeeuw*, on *07 Feb 2015 23:37:11 UTC*.\n",
	" - **Cloc: write errors to stderr instead of stdout** (#f5e25e4) by *Thomas de Zeeuw*, on *07 Feb 2015 23:36:41 UTC*.\n",
	" - **Cloc: clean path in every count function** (#d557a6c) by *Thomas de Zeeuw*, on *07 Feb 2015 14:23:13 UTC*.\n",
	" - **Cloc: stop after we encounter an error** (#c4ab7c7) by *Thomas de Zeeuw*, on *07 Feb 2015 14:05:31 UTC*.\n",
	" - **Add cloc** (#9014363) by *Thomas de Zeeuw*, on *07 Feb 2015 01:44:54 UTC*. <br/>\n*Cloc counts the number of lines of code in a given file or directory*.\n",
	" - **Add serve** (#a7a6823) by *Thomas de Zeeuw*, on *07 Feb 2015 01:44:33 UTC*. <br/>\n*Simply serve static files from a given directory on a given port*.\n",
	" - **init()** (#7539f40) by *Thomas de Zeeuw*, on *07 Feb 2015 01:43:58 UTC*.\n",
}

// Must match up with expectedCommits and expectedChangelog above.
var commitInput = []string{
	`hash: b6652b0
author: Thomas de Zeeuw
date: 2015-04-20T00:12:13+02:00
ref: HEAD -> master, tag: v0.2, origin/master
title: Serve: Allow first argument as directory
message: Now you can call ` + "`serve dir` or `server -d dir`" + `.

==============================`,

	`hash: 7c3924c
author: Thomas de Zeeuw
date: 2015-02-19T20:33:13+01:00
ref:
title: Move Gocover badges to the subdirectories
message:
==============================`,

	`hash: f7705fc
author: Thomas de Zeeuw
date: 2015-02-19T20:31:43+01:00
ref:
title: Drop godoc badge from main readme
message:
==============================`,

	`hash: 5f5e82f
author: Thomas de Zeeuw
date: 2015-02-19T20:27:07+01:00
ref:
title: all: adhere to 80 character line limit
message:
==============================`,

	`hash: 14600c3
author: Thomas de Zeeuw
date: 2015-02-19T20:25:33+01:00
ref:
title: cloc: skip file not found tests
message: same as previous commit, but these changes didn't get commited.

==============================`,

	`hash: 14d54f7
author: Thomas de Zeeuw
date: 2015-02-19T20:24:08+01:00
ref:
title: cll, cloc: skip file not found tests
message: The tests fail on travi-ci.org. The problem is likely that the error
message comes from the os and are therefor not platform independent.

==============================`,

	`hash: 143d198
author: Thomas de Zeeuw
date: 2015-02-19T20:02:08+01:00
ref: tag: v0.1
title: Fix Gocover badge
message:
==============================`,

	`hash: 0c05371
author: Thomas de Zeeuw
date: 2015-02-19T19:59:10+01:00
ref:
title: Add Travis, Gocover and Godoc badges
message:
==============================`,

	`hash: fc75e22
author: Thomas de Zeeuw
date: 2015-02-19T19:56:38+01:00
ref:
title: Fix travis config
message: .travis.yml requires spaces, not tabs.

==============================`,

	`hash: 99d43ae
author: Thomas de Zeeuw
date: 2015-02-19T19:52:55+01:00
ref:
title: Add editorconfig
message:
==============================`,

	`hash: 843f6e4
author: Thomas de Zeeuw
date: 2015-02-19T19:52:46+01:00
ref:
title: Add Travis-ci config
message: Adds testing for go 1.3, 1.4 and tip.

==============================`,

	`hash: 8188c6f
author: Thomas de Zeeuw
date: 2015-02-08T15:59:54+01:00
ref:
title: All: update examples in readme
message:
==============================`,

	`hash: 141942b
author: Thomas de Zeeuw
date: 2015-02-08T02:53:56+01:00
ref:
title: Cloc: improve readme example
message: ` + "`cloc my_file my_folder` is clearer then `cloc cloc.go _testdata`" + `

==============================`,

	`hash: 1bfdcf0
author: Thomas de Zeeuw
date: 2015-02-08T02:53:05+01:00
ref:
title: Add Cll
message: Cll checks if all lines have a length within the maximum allowed length.

==============================`,

	`hash: 545aa9d
author: Thomas de Zeeuw
date: 2015-02-08T01:29:41+01:00
ref:
title: All: update copyright notice in each source file
message:
==============================`,

	`hash: c6305e9
author: Thomas de Zeeuw
date: 2015-02-08T00:37:11+01:00
ref:
title: Serve: write errors to stderr instead of stdout
message:
==============================`,

	`hash: f5e25e4
author: Thomas de Zeeuw
date: 2015-02-08T00:36:41+01:00
ref:
title: Cloc: write errors to stderr instead of stdout
message:
==============================`,

	`hash: d557a6c
author: Thomas de Zeeuw
date: 2015-02-07T15:23:13+01:00
ref:
title: Cloc: clean path in every count function
message:
==============================`,

	`hash: c4ab7c7
author: Thomas de Zeeuw
date: 2015-02-07T15:05:31+01:00
ref:
title: Cloc: stop after we encounter an error
message:
==============================`,

	`hash: 9014363
author: Thomas de Zeeuw
date: 2015-02-07T02:44:54+01:00
ref:
title: Add cloc
message: Cloc counts the number of lines of code in a given file or directory

==============================`,

	`hash: a7a6823
author: Thomas de Zeeuw
date: 2015-02-07T02:44:33+01:00
ref:
title: Add serve
message: Simply serve static files from a given directory on a given port

==============================`,

	`hash: 7539f40
author: Thomas de Zeeuw
date: 2015-02-07T02:43:58+01:00
ref:
title: init()
message:
==============================`,
}
