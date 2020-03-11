// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests for some of the internal functions.

package main

import (
	"fmt"
	"testing"
)

// Helpers to save typing in the test cases.
type u []uint64
type uu [][]uint64

type SplitTest struct {
	name   string
	input  u
	output uu
	signed bool
}

var (
	m2          = uint64(2)
	m1          = uint64(1)
	m0          = uint64(0)
	m1Spanning0 = ^uint64(0)     // -1 when signed.
	m2Spanning0 = ^uint64(0) - 1 // -2 when signed.
)

var splitTests = []SplitTest{
	// No need for a test for the empty case; that's picked off before splitIntoRuns.
	{"single value", u{1}, uu{u{1}}, false},
	{"out of order signed", u{3, 2, 1}, uu{u{1, 2, 3}}, true},
	{"out of order unsigned", u{3, 2, 1}, uu{u{1, 2, 3}}, false},
	{"gap at the beginning", u{1, 33, 32, 31}, uu{u{1}, u{31, 32, 33}}, true},
	{"gap in the middle, mixed order", u{33, 7, 32, 31, 9, 8}, uu{u{7, 8, 9}, u{31, 32, 33}}, true},
	{"gaps throughout", u{33, 44, 1, 32, 45, 31}, uu{u{1}, u{31, 32, 33}, u{44, 45}}, true},
	{"values spanning 0 signed", u{m1, m0, m1Spanning0, m2, m2Spanning0}, uu{u{m2Spanning0, m1Spanning0, m0, m1, m2}}, true},
	{"values spanning 0 unsigned", u{m1, m0, m1Spanning0, m2, m2Spanning0}, uu{u{m0, m1, m2}, u{m2Spanning0, m1Spanning0}}, false},
}

func TestSplitIntoRuns(t *testing.T) {
	for _, test := range splitTests {
		t.Run(test.name, func(t *testing.T) {
			values := make([]Value, len(test.input))
			for i, v := range test.input {
				values[i] = Value{"", "", v, test.signed, fmt.Sprint(v)}
			}
			runs := splitIntoRuns(values)
			if len(runs) != len(test.output) {
				t.Fatalf("got %d runs; expected %d", len(runs), len(test.output))
			}
			for i, run := range runs {
				if len(run) != len(test.output[i]) {
					t.Fatalf("got %v; expected %v", runs, test.output)
				}
				for j, v := range run {
					if v.value != test.output[i][j] {
						t.Fatalf("got %v; expected %v", runs, test.output)
					}
				}
			}

		})
	}
}
