// Copyright (C) 2014 Thomas de Zeeuw.
// Licensed onder the MIT license that can be found in the LICENSE file.

package main

import "testing"

func TestGetTreeOrder(t *testing.T) {
	type test struct {
		input    map[string]Tree
		expected []string
	}

	tests := []test{
		{map[string]Tree{"c": {"c", false, nil}, "b": {"b", true, nil}, "a": {"a", false, nil}},
			[]string{"a", "c", "b"}},
	}

	for _, test := range tests {
		order := getTreeOrder(test.input)

		if len(order) != len(test.expected) {
			t.Errorf("Expected the order to have %d keys, but got %d", len(order), len(test.expected))
			continue
		}

		for i, o := range order {
			if o != test.expected[i] {
				t.Errorf("Expected %v, but got %v", test.expected, order)
			}
		}
	}
}
