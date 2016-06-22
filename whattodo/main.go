package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// For testing, see main_test.go init.
var printf = func(format string, a ...interface{}) (int, error) {
	return fmt.Printf(format, a...)
}

var exit = func(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

// Todo is a todo item with a message, a filepath and the line of the file.
type todo struct {
	text string
	path string
	line int
}

type ts []todo

func (ts ts) Len() int {
	return len(ts)
}

func (ts ts) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts ts) Less(i, j int) bool {
	if ts[i].path == ts[j].path {
		return ts[i].line < ts[j].line
	}
	return ts[i].path < ts[j].path
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [directory]", os.Args[0])
	}
	flag.Parse()

	path := filepath.Join("./", flag.Arg(0))
	if path == "" {
		path = "."
	}

	// Open the path and get stats from it. File is closed once we have the stats
	// if they return an error or not.
	f, err := os.Open(path)
	if err != nil {
		exit("Can't find path %s.", path)
	} else if s, err := f.Stat(); err != nil {
		exit("Can't get info from %s.", path)
	} else if !s.IsDir() {
		exit("Can only get todos from a whole directory.")
	}
	f.Close()

	// Create a new token file set and let Go do the heavy lifting.
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// For each package, each file within that package do something with the
	// comments in that file.
	var todos ts
	for _, pkg := range pkgs {
		for path, file := range pkg.Files {
			for _, cmt := range file.Comments {
				// Wether or not the last line had a todo item, used in adding the new
				// lines to the item (like this comment).
				var lastLineTodo bool

				for _, line := range strings.Split(cmt.Text(), "\n") {
					line = strings.TrimSpace(line)
					i := strings.Index(line, ":")

					if len(line) > 4 && strings.ToLower(line[:4]) == "todo" && i != -1 {
						// Drop anything before the first colon. This way we support
						// "todo:" and "todo (name):"
						item := strings.TrimSpace(line[i+1:])
						item = strings.ToUpper(item[:1]) + item[1:]

						// Recalculate the position of this line. First we need to get the
						// ofset of the current position. Then we add the position of this
						// line and convert it into a position again.
						// BUG: ast.CommentGroup.Text() reduces multiple empty lines to
						// one, which results in incorrect positions.
						f := fset.File(cmt.Pos())
						offset := f.Offset(cmt.Pos())
						offset += strings.Index(cmt.Text(), line) + len(line)
						pos := f.Pos(offset)

						// The position we can use in getting the line of the todo item.
						todos = append(todos, todo{item, path, fset.Position(pos).Line})
						lastLineTodo = true
					} else if lastLineTodo {
						// If the last line was a todo item, this line might be apart of
						// it. We'll break it on a empty line.
						if line == "" {
							lastLineTodo = false
							continue
						}

						// Add the line text to the todo text.
						todos[len(todos)-1].text += " " + line
					}
				}
			}
		}
	}

	// Wel done.
	if len(todos) == 0 {
		printf("Nothing, you have done all you needed to do (atleast according " +
			"to the lack of todos in the source code).")
		return
	}

	// Print the found todo items.
	sort.Sort(todos)
	for _, item := range todos {
		printf("  - [ ] %s (%s, line %d).\n", item.text, item.path, item.line)
	}
}
