// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Marshaling test: test marshaler intefaces.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type MarshalText int

const (
	Monday MarshalText = iota
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
	ckUnmarshal(`{"Monday":"test"}`, map[MarshalText]string{Monday: "test"})
	ckUnmarshal(`{"Tuesday":"test"}`, map[MarshalText]string{Tuesday: "test"})
	ckUnmarshal(`{"Wednesday":"test"}`, map[MarshalText]string{Wednesday: "test"})
	ckUnmarshal(`{"Thursday":"test"}`, map[MarshalText]string{Thursday: "test"})
	ckUnmarshal(`{"Friday":"test"}`, map[MarshalText]string{Friday: "test"})
	ckUnmarshal(`{"Saturday":"test"}`, map[MarshalText]string{Saturday: "test"})
	ckUnmarshal(`{"Sunday":"test"}`, map[MarshalText]string{Sunday: "test"})
	ckUnmarshal(`{"Christmas":"test"}`, nil)
}

const panicPrefix = "marshal_text.go: "

func ckNoError(err error) {
	if err != nil {
		panic(panicPrefix + err.Error())
	}
}

func ckMarshal(e MarshalText) {
	expected := `{"` + e.String() + `":"test"}`
	m := map[MarshalText]string{e: "test"}
	raw, err := json.Marshal(m)
	ckNoError(err)
	compactBuff := bytes.Buffer{}
	err = json.Compact(&compactBuff, raw)
	ckNoError(err)
	got := string(compactBuff.Bytes())
	if got != expected {
		panic(fmt.Sprintf("%s text.Marshal got '%s', expected '%s'", panicPrefix, got, expected))
	}
}

func ckUnmarshal(raw string, expected map[MarshalText]string) {
	failed := false
	got := make(map[MarshalText]string)
	err := json.Unmarshal([]byte(raw), &got)
	if expected == nil {
		if err == nil {
			panic(panicPrefix + "expected error")
		}
	} else {
		ckNoError(err)
		for k, v := range expected {
			gotV, ok := got[k]
			if !ok || gotV != v {
				failed = true
				break
			}
		}

		if failed {
			panic(fmt.Sprintf("%s text.Unmarshal got '%#v', expected '%#v'", panicPrefix, got, expected))
		}
	}
}
