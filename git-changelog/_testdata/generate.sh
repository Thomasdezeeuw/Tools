#!/bin/bash

# Generate generates some new test data from the current git repo. It will
# create a new file for the log and not overwrite any previous testdata.

format="hash: %h
author: %an
date: %cI
ref: %D
title: %s
message: %b
==============================%n"

for (( i = 1; i <= 100; i++ )); do
	file="git.$i.log"
	if [ ! -e "$file" ] ; then
		git log --format="$format" > "$file"
		break
	fi
done
