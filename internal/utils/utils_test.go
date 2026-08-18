// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		a     string
		b     []byte
		equal bool
	}{
		{"", nil, true},
		{"", []byte{}, true},
		{" ", []byte{' '}, true},
		{"hello", []byte("world"), false},
		{"hello", []byte("hello"), true},
	}

	for _, test := range tests {
		equal := Equal(test.a, test.b)

		if test.equal && !equal {
			t.Errorf("expected %s equal %s", test.a, test.b)
		} else if !test.equal && equal {
			t.Errorf("expected %s not equal %s", test.a, test.b)
		}
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		body     string
		expected int
		ok       bool
	}{
		{"", 0, false},
		{"0", 0, true},
		{"1234", 1234, true},
		{"-1234", -1234, true},
		{"hello", 0, false},
	}

	for _, test := range tests {
		t.Run(test.body, func(t *testing.T) {
			target, ok := Int[int](test.body)

			if test.ok && !ok {
				t.Fatal("expected ok")
			} else if !test.ok && ok {
				t.Fatal("expected not ok")
			}

			if test.ok {
				if test.expected != target {
					t.Errorf("expected %d got %d", test.expected, target)
				}
			}
		})
	}
}

func TestFields(t *testing.T) {
	tests := []struct {
		desc     string
		body     string
		expected []string
	}{
		{"Empty", "", []string{}},
		{"NoFields", "hello", []string{"hello"}},
		{"Field", "hello world", []string{"hello", "world"}},
		{"Leading", "   hello world", []string{"hello", "world"}},
		{"Trailing", "hello world   ", []string{"hello", "world"}},
		{"ManySpaces", "   hello   world   ", []string{"hello", "world"}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			target := []string{}

			for _, t := range Fields(test.body) {
				target = append(target, t)
			}

			if !cmp.Equal(test.expected, target) {
				t.Error(cmp.Diff(test.expected, target))
			}
		})
	}
}

func TestLines(t *testing.T) {
	tests := []struct {
		desc     string
		body     string
		expected []string
	}{
		{"Empty", "", []string{}},
		{"NoLines", "hello", []string{"hello"}},
		{"Line", "hello\nworld", []string{"hello", "world"}},
		{"Leading", "\nhello\nworld", []string{"hello", "world"}},
		{"Trailing", "hello\nworld\n", []string{"hello", "world"}},
		{"ManyLines", "\n\n\nhello\n\n\nworld\n\n\n", []string{"hello", "world"}},
		{"Spaces", "\n\n\n   hello   \n\n\n   world   \n\n\n", []string{"   hello   ", "   world   "}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			target := []string{}

			for _, t := range Lines(test.body) {
				target = append(target, t)
			}

			if !cmp.Equal(test.expected, target) {
				t.Error(cmp.Diff(test.expected, target))
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		body   []byte
		prefix string
		has    bool
	}{
		{nil, "", true},
		{[]byte{}, "", true},
		{[]byte(" "), " ", true},
		{[]byte("hello"), "hello", true},
		{[]byte("hello world"), "hello", true},
		{[]byte("foo bar"), "hello", false},
	}

	for _, test := range tests {
		has := HasPrefix(test.body, test.prefix)

		if test.has && !has {
			t.Errorf("expected %s has prefix %s", test.body, test.prefix)
		} else if !test.has && has {
			t.Errorf("expected %s not have prefix %s", test.body, test.prefix)
		}
	}
}

func TestIndexByte(t *testing.T) {
	tests := []struct {
		desc     string
		body     string
		expected int
	}{
		{"Empty", "", -1},
		{"NoByte", "hello", -1},
		{"Byte", "hello world", 5},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			target := IndexByte(test.body, isSpace)

			if test.expected != target {
				t.Errorf("expected %d got %d", test.expected, target)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		desc   string
		body   string
		before string
		after  string
		split  bool
	}{
		{"Empty", "", "", "", false},
		{"NoSplit", "hello", "hello", "", false},
		{"Split", "hello world", "hello", "world", true},
		{"Delims", "foo bar baz", "foo", "bar baz", true},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			before, after, split := Split(test.body, isSpace)

			if test.split && !split {
				t.Fatal("expected split")
			} else if !test.split && split {
				t.Fatal("unexpected split")
			}

			if test.split {
				if test.before != before {
					t.Errorf("expected before:\n%s\ngot:\n%s", test.before, before)
				}

				if test.after != after {
					t.Errorf("expected after:\n%s\ngot:\n%s", test.after, after)
				}

			} else {
				if test.body != before {
					t.Errorf("expected body==before:\n%s\ngot:\n%s", test.body, before)
				}
			}
		})
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		desc     string
		body     string
		expected string
	}{
		{"Empty", "", ""},
		{"NoSpace", "hello", "hello"},
		{"Leading", "   hello", "hello"},
		{"Trailing", "hello   ", "hello"},
		{"Leading/Trailing", "   hello   ", "hello"},
		{"OnlySpace", "   ", ""},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			target := TrimSpace(test.body)

			if test.expected != target {
				t.Errorf("expected:\n%s\ngot:\n%s", test.expected, target)
			}
		})
	}
}

func TestTrimLeft(t *testing.T) {
	tests := []struct {
		desc     string
		body     string
		expected string
	}{
		{"Empty", "", ""},
		{"NoSpace", "hello", "hello"},
		{"Leading", "   hello", "hello"},
		{"Trailing", "hello   ", "hello   "},
		{"OnlySpace", "   ", ""},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			target := TrimLeft(test.body)

			if test.expected != target {
				t.Errorf("expected:\n%s\ngot:\n%s", test.expected, target)
			}
		})
	}
}

func TestTrimRight(t *testing.T) {
	tests := []struct {
		desc     string
		body     string
		expected string
	}{
		{"Empty", "", ""},
		{"NoSpace", "hello", "hello"},
		{"Trailing", "hello   ", "hello"},
		{"Leading", "   hello", "   hello"},
		{"OnlySpace", "   ", ""},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			target := TrimRight(test.body)

			if test.expected != target {
				t.Errorf("expected:\n%s\ngot:\n%s", test.expected, target)
			}
		})
	}
}

func TestUint(t *testing.T) {
	tests := []struct {
		body     string
		expected uint
		ok       bool
	}{
		{"", 0, false},
		{"0", 0, true},
		{"1234", 1234, true},
		{"-1234", 0, false},
		{"hello", 0, false},
	}

	for _, test := range tests {
		t.Run(test.body, func(t *testing.T) {
			target, ok := Uint[uint](test.body)

			if test.ok && !ok {
				t.Fatal("expected ok")
			} else if !test.ok && ok {
				t.Fatal("expected not ok")
			}

			if test.ok {
				if test.expected != target {
					t.Errorf("expected %d got %d", test.expected, target)
				}
			}
		})
	}
}
