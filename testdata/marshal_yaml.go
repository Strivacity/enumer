// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Marshaling test: test marshaler intefaces.

package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type MarshalYaml int

const (
	Monday MarshalYaml = iota
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
	ckUnmarshal(`Monday`, Monday)
	ckUnmarshal(`Tuesday`, Tuesday)
	ckUnmarshal(`Wednesday`, Wednesday)
	ckUnmarshal(`Thursday`, Thursday)
	ckUnmarshal(`Friday`, Friday)
	ckUnmarshal(`Saturday`, Saturday)
	ckUnmarshal(`Sunday`, Sunday)
	ckUnmarshal(`Christmas`, 127)

}

const panicPrefix = "marshal_yaml.go: "

func ckNoError(err error) {
	if err != nil {
		panic(panicPrefix + err.Error())
	}
}

func ckMarshal(e MarshalYaml) {
	raw, err := yaml.Marshal(&e)
	ckNoError(err)
	got := string(raw)
	expected := e.String() + "\n"
	if got != expected {
		panic(fmt.Sprintf("%s yaml.Marshal got '%s', expected '%s'", panicPrefix, got, expected))
	}
}

func ckUnmarshal(raw string, expected MarshalYaml) {
	var got MarshalYaml = 0
	err := yaml.Unmarshal([]byte(raw), &got)
	if !expected.IsAMarshalYaml() {
		if err == nil {
			panic(panicPrefix + "expected error")
		}
	} else {
		ckNoError(err)
		if got != expected {
			panic(fmt.Sprintf("%s yaml.Unmarshal got '%s', expected '%s'", panicPrefix, got, expected))
		}
	}
}
