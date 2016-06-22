// Copyright (C) 2014 Thomas de Zeeuw.
// Licensed onder the MIT license that can be found in the LICENSE file.

// Tree creates a tree like representation of a given directory.
//
// Create a tree from the current directory:
// 	$ tree
// 	tree
// 	└─ main.go
// 	└─ main_test.go
// 	└─ test_files
// 	 └─ file1
// 	 └─ test_files2
// 	  └─ file2
// 	 └─ test_files3
// 	  └─ file3
//
// Creating a tree of a given directory:
// 	$ tree test_files
// 	test_files
// 	└─ file1
// 	└─ test_files2
// 	 └─ file2
// 	└─ test_files3
// 	 └─ file3
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// Tree is a small representation a file/directory
type Tree struct {
	Name      string          // Name of the file
	IsDir     bool            // Wether or not it's an directory
	Childeren map[string]Tree // Possible childeren (in case of a directory)
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		path = "."
	}

	if len(os.Args) >= 2 {
		if len(os.Args) != 2 {
			fmt.Print("Only accepting a single argument, the path to create the tree for.")
			return
		}

		path = os.Args[1]
	}

	trees, err := createTree(path)
	if err != nil {
		fmt.Print(err)
		return
	}

	printTrees(path, trees, "")
}

func createTree(path string) (map[string]Tree, error) {
	path = filepath.Clean(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot open file %s.", path)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Cannot stat open file %s.", path)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("Error closing file %s.", path)
	}

	if !stat.Mode().IsDir() {
		return nil, fmt.Errorf("%s is not a directory.", path)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	trees := map[string]Tree{}

	for _, file := range files {
		tree := Tree{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		}

		if tree.IsDir {
			childeren, err := createTree(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}

			tree.Childeren = childeren
		}

		trees[tree.Name] = tree
	}

	return trees, nil
}

// printTrees print the tree nicely formatted
func printTrees(path string, trees map[string]Tree, space string) {
	fmt.Printf("%s\n", filepath.Base(path))

	for _, k := range getTreeOrder(trees) {
		tree := trees[k]

		fmt.Printf("%s└─ ", space)

		if tree.IsDir {
			printTrees(tree.Name, tree.Childeren, space+" ")
		} else {
			fmt.Printf("%s\n", tree.Name)
		}
	}
}

// getTreeOrder create an order for a tree. Files first then directories.
func getTreeOrder(trees map[string]Tree) []string {
	var files = make([]string, 0, len(trees))
	var dirs = make([]string, 0, len(trees))

	for key, tree := range trees {
		if tree.IsDir {
			dirs = append(dirs, key)
		} else {
			files = append(files, key)
		}
	}

	sort.Strings(files)
	sort.Strings(dirs)

	return append(files, dirs...)
}
