// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains simple golden tests for various examples.
// Besides validating the results when the implementation changes,
// it provides a way to look at the generated code without having
// to execute the print statements in one's head.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Golden represents a test case.
type Golden struct {
	name  string
	input string // input; the package clause is provided when running the test.
	cfg   config
}

var tests = []Golden{
	{"day", dayIn, config{}},
	{"offset", offsetIn, config{}},
	{"gap", gapIn, config{}},
	{"num", numIn, config{}},
	{"unum", unumIn, config{}},
	{"prime", primeIn, config{}},
	{"primeJson", primeIn, config{json: true}},
	{"primeText", primeIn, config{text: true}},
	{"primeYaml", primeIn, config{yaml: true}},
	{"primeSql", primeIn, config{sql: true}},
	{"primeJsonAndSql", primeIn, config{json: true, sql: true}},
	{"dayTrimPrefix", dayPrefixIn, config{trimPrefix: "Day"}},
	{"dayTrimPrefixMultiple", dayPrefixMultipleIn, config{trimPrefix: "Day,Night"}},
	{"dayWithPrefix", dayIn, config{addPrefix: "Day"}},
	{"dayTrimAndPrefix", dayPrefixIn, config{trimPrefix: "Day", addPrefix: "Night"}},
	{"daySet", dayIn, config{setDelimiter: ","}},
	{"daySetTrimPrefix", dayPrefixIn, config{setDelimiter: ",", trimPrefix: "Day"}},
	{"daySetStrict", dayIn, config{setDelimiter: ",", strictSet: true}},
}

// Each example starts with "type XXX [u]int", with a single space separating them.

// Simple test: enumeration of type int starting at 0.
const dayIn = `type Day int
const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
`

const dayPrefixIn = `type Day int
const (
	DayMonday Day = iota
	DayTuesday
	DayWednesday
	DayThursday
	DayFriday
	DaySaturday
	DaySunday
)
`

const dayPrefixMultipleIn = `type Day int
const (
	DayMonday Day = iota
	NightTuesday
	DayWednesday
	NightThursday
	DayFriday
	NightSaturday
	DaySunday
)
`

// Enumeration with an offset.
// Also includes a duplicate.
const offsetIn = `type Number int
const (
	_ Number = iota
	One
	Two
	Three
	AnotherOne = One  // Duplicate; note that AnotherOne doesn't appear below.
)
`

// Gaps and an offset.
const gapIn = `type Gap int
const (
	Two Gap = 2
	Three Gap = 3
	Five Gap = 5
	Six Gap = 6
	Seven Gap = 7
	Eight Gap = 8
	Nine Gap = 9
	Eleven Gap = 11
)
`

// Signed integers spanning zero.
const numIn = `type Num int
const (
	m_2 Num = -2 + iota
	m_1
	m0
	m1
	m2
)
`

// Unsigned integers spanning zero.
const unumIn = `type Unum uint
const (
	m_2 Unum = iota + 253
	m_1
)

const (
	m0 Unum = iota
	m1
	m2
)
`

// Enough gaps to trigger a map implementation of the method.
// Also includes a duplicate to test that it doesn't cause problems
const primeIn = `type Prime int
const (
	p2 Prime = 2
	p3 Prime = 3
	p5 Prime = 5
	p7 Prime = 7
	p77 Prime = 7 // Duplicate; note that p77 doesn't appear below.
	p11 Prime = 11
	p13 Prime = 13
	p17 Prime = 17
	p19 Prime = 19
	p23 Prime = 23
	p29 Prime = 29
	p37 Prime = 31
	p41 Prime = 41
	p43 Prime = 43
)
`

func TestGolden(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runGoldenTest(t, test)
		})
	}
}

func runGoldenTest(t *testing.T, test Golden) {
	goldenFile := fmt.Sprintf("./goldendata/%v.golden", test.name)
	expectedBytes, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Error(err)
	}
	expected := string(expectedBytes)

	var g Generator
	file := test.name + ".go"
	input := "package test\n" + test.input

	dir, err := ioutil.TempDir("", "stringer")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()

	absFile := filepath.Join(dir, file)
	err = ioutil.WriteFile(absFile, []byte(input), 0644)
	if err != nil {
		t.Error(err)
	}
	g.parsePackage([]string{absFile}, nil)
	// Extract the name and type of the constant from the first line.
	tokens := strings.SplitN(test.input, " ", 3)
	if len(tokens) != 3 {
		t.Fatal("need type declaration on first line")
	}
	g.generate(tokens[1], test.cfg)
	got := string(g.format())
	if got != expected {
		// Use this to help build a golden text when changes are needed
		//err = ioutil.WriteFile(goldenFile, []byte(got), 0644)
		//if err != nil {
		//	t.Error(err)
		//}
		t.Errorf("got\n====\n%s====\nexpected\n====\n%s====\n", got, expected)
	}
}
