// Copyright (C) 2015 Thomas de Zeeuw.
//
// Licensed onder the MIT license that can be found in the LICENSE file.

// TODO: add options for the depth, so not all underlying files and directories
// are checked.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type longLine struct {
	File       string
	LineNumber int
	Length     int
}

// The maximum length of a line. It defaults to 80, but get overwritten by the
// -l and length flags.
var maxLength = 80

// The length of a tab, \t is counted as a single character, but in most
// editors it's 2 or 4 spaces.
var tabLength = 2

// Descriptions used for the flags.
const (
	maxLengthDesc = "Maximum allowed length of a line, defaults to 80"
	tabLengthDesc = "The length of a tab, defaults to 2"
)

func init() {
	flag.IntVar(&maxLength, "l", maxLength, maxLengthDesc)
	flag.IntVar(&maxLength, "length", maxLength, maxLengthDesc)
	flag.IntVar(&tabLength, "t", tabLength, tabLengthDesc)
	flag.IntVar(&tabLength, "tab", tabLength, tabLengthDesc)
}

func main() {
	flag.Parse()
	files := flag.Args()

	// Default to counting the code in the current diretory.
	if len(files) == 0 {
		files = []string{"./"}
	}

	for _, path := range files {
		longLines, err := check(path, maxLength)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			return
		}

		// Report our results to the user.
		for _, longLine := range longLines {
			fmt.Printf("%s line %d is %d characters long.\n", longLine.File,
				longLine.LineNumber, longLine.Length)
		}
	}
}

// Check checks if all lines in a directory or files are shorter then the
// maximum allowed length.
func check(path string, maxLength int) ([]longLine, error) {
	path = filepath.Clean(path)

	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot open file %s.", path)
	}

	// Get the file information.
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Cannot stat open file %s.", path)
	}

	// Close the file, we won't need it anymore.
	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("Error closing file %s.", path)
	}

	// Count the number of lines in the file or directory.
	if stat.Mode().IsDir() {
		return checkDir(path, maxLength)
	} else {
		return checkFile(path, maxLength)
	}
}

// CheckDir checks for all files in a directory if they have lines with a
// length that surpasses the maximum allowed length.
func checkDir(dirpath string, maxLength int) ([]longLine, error) {
	dirpath = filepath.Clean(dirpath)

	// Grap all the files and directories in the given directory.
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	// A counter for the number of files so we can make sure we wait for every
	// one of them.
	fileCounter := len(files)

	// A channel for the long lines and one for possible errors.
	resultChannel := make(chan []longLine, fileCounter)
	errorChannel := make(chan error, 1)

	// For each file/directory in the directory.
	for _, file := range files {
		go func(path string) {
			path = filepath.Join(dirpath, path)

			// Check if the directory or file has long lines
			longLines, err := check(path, maxLength)
			if err != nil {
				errorChannel <- err
				return
			}

			resultChannel <- longLines
		}(file.Name())
	}

	var allLongLines []longLine

	// Wait for a response from each file/directory and either respond with an
	// error or add the count to the total number of lines.
	for fileCounter > 0 {
		select {
		case err := <-errorChannel:
			return nil, err
		case longLines := <-resultChannel:
			allLongLines = append(allLongLines, longLines...)
		}

		// Need to wait for yet another less response.
		fileCounter--
	}

	// Cleanup.
	close(resultChannel)
	close(errorChannel)

	return allLongLines, nil
}

// CheckFile checks if a file has lines with a length that surpasses the
// maximum allowed length.
func checkFile(path string, maxLength int) ([]longLine, error) {
	path = filepath.Clean(path)

	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Create a reader that returns a single line at a time.
	scanner := bufio.NewScanner(file)

	// A line number, so we can provide which line is to long.
	lineNumber := 1

	var longLines []longLine
	for scanner.Scan() {
		// Get the actual line
		line := scanner.Text()

		// Replace tabs with a number of spaces
		line = strings.Replace(line, "\t", strings.Repeat(" ", tabLength), -1)

		// Check if the length of the line is longer that our allowed length.
		if len(line) > maxLength {
			longLines = append(longLines, longLine{path, lineNumber, len(line)})
		}

		lineNumber++
	}

	// Check for possible errors.
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return longLines, nil
}
