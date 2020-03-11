// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Set test: strict parsing.

package main

type StrictDay int

const (
	Monday StrictDay = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	ckStrictDaySetString(StrictDaySetSunday, "Sunday")
	ckStrictDaySetString(StrictDaySetSunday, "sunday")
	ckStrictDaySetString(StrictDaySetSunday|StrictDaySetMonday, "Monday Sunday")
	ckStrictDaySetStringError("Christmas")
	ckStrictDaySetStringError("Monday Christmas")
}

const panicPrefix = "set_strict.go: "

func ckStrictDaySetString(day StrictDaySet, str string) {
	d, err := StrictDaySetString(str)
	if err != nil {
		panic(panicPrefix + err.Error())
	}
	if d != day {
		panic(panicPrefix + str)
	}
}

func ckStrictDaySetStringError(str string) {
	d, err := StrictDaySetString(str)
	if err == nil {
		panic(panicPrefix + "expected error: " + str)
	}
	if d != 0 {
		panic(panicPrefix + "expected zero value: " + str)
	}
}
