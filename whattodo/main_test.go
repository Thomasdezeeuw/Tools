package main

import (
	"fmt"
	"os"
	"sort"
	"testing"
)

var (
	result  = ""
	oldArgs = os.Args
)

func init() {
	printf = func(format string, a ...interface{}) (int, error) {
		result += fmt.Sprintf(format, a...)
		return 0, nil
	}

	exit = func(format string, a ...interface{}) {
		panic(fmt.Sprintf(format, a...))
	}
}

func TestMain(t *testing.T) {
	defer reset()
	os.Args = []string{"", ""}
	main()

	expected := "Nothing, you have done all you needed to do (atleast " +
		"according to the lack of todos in the source code)."

	if result != expected {
		t.Fatalf("The develop needs to get to work: %s", result[7:])
	}
}

func TestMainTestData(t *testing.T) {
	defer reset()
	os.Args = []string{"", "testdata"}
	main()

	expected := `  - [ ] More testdata! (testdata\more_testdata.go, line 4).
  - [ ] Not enough o's we need to add more o's! (testdata\more_testdata.go, line 10).
  - [ ] Be more excited (testdata\testdata.go, line 5).
`

	if result != expected {
		t.Fatalf("Expected the result to be: \n%s \nBut got: \n%s", expected, result)
	}
}

func TestTodosSort(t *testing.T) {
	todos := ts{
		{"Msg", "Path2", 1},
		{"Msg", "Path1", 2},
		{"Msg", "Path1", 1},
	}
	sort.Sort(todos)

	expected := []todo{
		{"Msg", "Path1", 1},
		{"Msg", "Path1", 2},
		{"Msg", "Path2", 1},
	}

	for i, todo := range todos {
		if todo != expected[i] {
			t.Fatalf("Expected the todos to be sort like %+v, got %+v",
				expected, todos)
		}
	}
}

func reset() {
	os.Args = oldArgs
	result = ""
}
