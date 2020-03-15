// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Marshaling test: test marshaler intefaces.

package main

import (
	"encoding/json"
	"fmt"
)

type MarshalJSON int

const (
	Monday MarshalJSON = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	ckMarshal(Monday)
	ckMarshal(Tuesday)
	ckMarshal(Wednesday)
	ckMarshal(Thursday)
	ckMarshal(Friday)
	ckMarshal(Saturday)
	ckMarshal(Sunday)
	ckMarshal(127)
	ckMarshal(-127)
	ckUnmarshal(`"Monday"`, Monday)
	ckUnmarshal(`"Tuesday"`, Tuesday)
	ckUnmarshal(`"Wednesday"`, Wednesday)
	ckUnmarshal(`"Thursday"`, Thursday)
	ckUnmarshal(`"Friday"`, Friday)
	ckUnmarshal(`"Saturday"`, Saturday)
	ckUnmarshal(`"Sunday"`, Sunday)
	ckUnmarshal(`"Christmas"`, 128)
}

const panicPrefix = "marshal_json.go: "

func ckNoError(err error) {
	if err != nil {
		panic(panicPrefix + err.Error())
	}
}

func ckMarshal(e MarshalJSON) {
	raw, err := json.Marshal(e)
	ckNoError(err)
	got := string(raw)
	expected := `"` + e.String() + `"`
	if got != expected {
		panic(fmt.Sprintf("%s json.Marshal got '%s', expected '%s'", panicPrefix, got, expected))
	}
}

func ckUnmarshal(raw string, expected MarshalJSON) {
	var got MarshalJSON = 0
	err := json.Unmarshal([]byte(raw), &got)
	if !expected.IsAMarshalJSON() {
		if err == nil {
			panic(panicPrefix + "expected error")
		}
	} else {
		ckNoError(err)
		if got != expected {
			panic(fmt.Sprintf("%s json.Unmarshal got '%s', expected '%s'", panicPrefix, got, expected))
		}
	}
}
