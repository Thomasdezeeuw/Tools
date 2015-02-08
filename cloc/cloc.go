// Copyright (C) 2015 Thomas de Zeeuw.
//
// Licensed onder the MIT license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	newLineBytes  = []byte("\n")
	doubleNewLine = regexp.MustCompile(`\n\s*\n`)
)

type language struct {
	OneLine    []*regexp.Regexp // Regexp to detect single line comments.
	MultiLine  []*regexp.Regexp // Regexp to detect multi line comments.
	Extentions []string         // Known file extentions for the language.
}

func main() {
	totalCount := 0
	files := getFileOptions(os.Args)

	for _, path := range files {
		count, err := count(path)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			return
		}

		totalCount += count

		// Current directory is clearer then ./.
		if path == "./" {
			path = "Current directory"
		}

		// Inform the user about the number of lines.
		fmt.Printf("%s contains %d lines of code.\n", path, count)
	}

	// Print a total of number of lines if the user requested info on more then
	// one file or directory.
	if len(files) > 1 {
		fmt.Printf("Total number of lines: %d.\n", totalCount)
	}
}

// GetFileOptions gets the files we need to count from the command line
// arguments.
func getFileOptions(args []string) []string {
	var files []string
	if len(args) > 1 {
		skip := false

		for _, arg := range args[1:] {
			// If the last argument started with - or --, we'll skip this argument
			// aswell.
			if skip {
				skip = false
				continue
			}

			// If the argument starts with - or --, we need to skip it aswell as the
			// next one.
			if (len(arg) >= 1 && arg[:1] == "-") || (len(arg) >= 2 && arg[:2] == "--") {
				skip = true
				continue
			}

			arg = strings.TrimSpace(arg)
			if arg == "" {
				continue
			}

			// Otherwise we will count the number of lines in the file or directory.
			files = append(files, filepath.Clean(arg))
		}
	}

	// Default to counting the code in the current diretory.
	if len(files) == 0 {
		files = []string{"./"}
	}

	return files
}

// Count counts the number of lines of code in a file or all files in a
// directory. Path can either be a file or a directory, in case of a directory
// all subdirectories will be counted aswell.
//
// If a file in not detected as a source file it will return 0 but without an
// error.
//
// Possible returned errors are mostly related to not being able to open or
// read the given path.
func count(path string) (int, error) {
	path = filepath.Clean(path)

	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("Cannot open file %s.", path)
	}

	// Get the file information.
	stat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("Cannot stat open file %s.", path)
	}

	// Close the file, we won't need it anymore.
	err = file.Close()
	if err != nil {
		return 0, fmt.Errorf("Error closing file %s.", path)
	}

	// Count the number of lines in the file or directory.
	if stat.Mode().IsDir() {
		return countDir(path)
	} else {
		return countFile(path)
	}
}

// CountDir counts the number of lines of code in all files in a given
// directory and it's subdirectories.
//
// TODO(Thomas): support for ignoring directories.
// TODO(Thomas): support for setting recursive level.
func countDir(dirpath string) (int, error) {
	dirpath = filepath.Clean(dirpath)

	// Grap all the files and directories in the given directory.
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return 0, err
	}

	// A counter for the counted number of lines and one for the number of files
	// so we can make sure we wait for every one of them.
	codeLineCounter := 0
	fileCounter := len(files)

	// A channel for the number of counted files and one for possible errors.
	countChannel := make(chan int, fileCounter)
	errorChannel := make(chan error, 1)

	// For each file/directory in the directory.
	for _, file := range files {
		go func(path string) {
			path = filepath.Join(dirpath, path)

			// Count the number of lines of the directory or file.
			count, err := count(path)
			if err != nil {
				errorChannel <- err
				return
			}

			countChannel <- count
		}(file.Name())
	}

	// Wait for a response from each file/directory and either respond with an
	// error or add the count to the total number of lines.
	for fileCounter > 0 {
		select {
		case count := <-countChannel:
			codeLineCounter += count
		case err := <-errorChannel:
			return 0, err
		}

		// Need to wait for yet another less response.
		fileCounter--
	}

	// Cleanup.
	close(countChannel)
	close(errorChannel)

	return codeLineCounter, nil
}

// CountFile counts the number of lines of code in a single given file, if the
// file is not a source file then we'll return 0, but not an error.
//
// BUG(Thomas): Doesn't work with all encodings, BOM enconding generally
// doesn't work.
func countFile(path string) (int, error) {
	path = filepath.Clean(path)

	// Get the langauge from the file path.
	lang := getLanguage(path)

	// Not a source file so we'll return a count of 0.
	if lang == languages["unkown"] {
		return 0, nil
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	// If the language supports single line comments remove them.
	if lang.OneLine != nil {
		// Replace all single line comments with a empty new line.
		for _, oneLine := range lang.OneLine {
			fileBytes = oneLine.ReplaceAll(fileBytes, newLineBytes)
		}
	}

	// If the language supports multi line comments remove them.
	if lang.MultiLine != nil {
		// Replace all multi line comments with nothing.
		for _, multiLine := range lang.MultiLine {
			fileBytes = multiLine.ReplaceAll(fileBytes, []byte{})
		}
	}

	// Replace every double new line with a single new line, basicly dropping
	// empty lines.
	fileBytes = doubleNewLine.ReplaceAll(fileBytes, newLineBytes)

	// Trim null bytes and space.
	fileBytes = bytes.Trim(fileBytes, "\x00")
	fileBytes = bytes.TrimSpace(fileBytes)

	// Only count the number of lines if the file is not empty.
	count := 0
	if len(fileBytes) != 0 {
		// The number of line endings plus 1 is the number of lines.
		count = bytes.Count(fileBytes, newLineBytes) + 1
	}

	return count, nil
}

// GetLanguage detects the language based on the file extention.
func getLanguage(path string) *language {
	// Get the extention from the path.
	ext := strings.TrimPrefix(filepath.Ext(path), ".")

	// Check if it matches a known extention of one of the languages.
	for _, lang := range languages {
		for _, langExt := range lang.Extentions {
			if langExt == ext {
				return lang
			}
		}
	}

	// Don't known the language, so we'll return the unkown language.
	return languages["unkown"]
}
