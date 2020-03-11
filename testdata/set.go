// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Set test: lax parsing.

package main

import "fmt"

type LaxDay int

const (
	Monday LaxDay = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	ck(LaxDaySetMonday, "Monday")
	ck(LaxDaySetTuesday, "Tuesday")
	ck(LaxDaySetWednesday, "Wednesday")
	ck(LaxDaySetThursday, "Thursday")
	ck(LaxDaySetFriday, "Friday")
	ck(LaxDaySetSaturday, "Saturday")
	ck(LaxDaySetSunday, "Sunday")
	ck(LaxDaySetMonday|LaxDaySetTuesday, "Monday Tuesday")
	ck(LaxDaySetTuesday|LaxDaySetMonday, "Monday Tuesday")
	ck(127, "Monday Tuesday Wednesday Thursday Friday Saturday Sunday")
	ckDaySetElements(LaxDaySetSunday, []LaxDay{Sunday})
	ckDaySetElements(LaxDaySetSunday|LaxDaySetMonday, []LaxDay{Sunday, Monday})
	ckDaySetTest("Saturday Sunday", "Monday", false)
	ckDaySetTest("Saturday Sunday", "Monday Tuesday Wednesday Thursday Friday", false)
	ckDaySetTest("Saturday Sunday", "Monday Sunday", false)
	ckDaySetTest("Saturday Sunday", "Saturday", true)
	ckDaySetTest("Saturday Sunday", "Sunday", true)
	ckDaySetTest("Saturday Sunday", "Saturday Sunday", true)
	ckLaxDaySetString(LaxDaySetSunday, "Sunday")
	ckLaxDaySetString(LaxDaySetSunday, "sunday")
	ckLaxDaySetString(LaxDaySetSunday|LaxDaySetMonday, "Monday Sunday")
	ckLaxDaySetString(LaxDaySetMonday, "Monday Christmas")
}

const panicPrefix = "set.go: "

func ck(daySet LaxDaySet, str string) {
	if fmt.Sprint(daySet) != str {
		panic(fmt.Sprintf("%s got '%s', expected '%s'", panicPrefix, fmt.Sprint(daySet), str))
	}
}

func ckDaySetElements(daySet LaxDaySet, expectedElems []LaxDay) {
	elems := daySet.Elements()
	if len(elems) != len(expectedElems) {
		panic(fmt.Sprintf("%sexpected %d elements, got %d", panicPrefix, len(expectedElems), len(elems)))
	}
OUTER:
	for _, expected := range expectedElems {
		for _, got := range elems {
			if got == expected {
				continue OUTER
			}
		}
		panic(fmt.Sprintf("%smissing element %s", panicPrefix, expected))
	}
}

func ckDaySetTest(ds string, m string, expected bool) {
	daySet, err := LaxDaySetString(ds)
	ckNoError(err)
	mask, err := LaxDaySetString(m)
	ckNoError(err)
	if daySet.Test(mask) != expected {
		panic(fmt.Sprintf("%stesting '%s', on '%s' is not %s", panicPrefix, mask, daySet, expected))
	}
}

func ckLaxDaySetString(daySet LaxDaySet, str string) {
	d, err := LaxDaySetString(str)
	ckNoError(err)
	if d != daySet {
		panic(panicPrefix + str)
	}
}

func ckNoError(err error) {
	if err != nil {
		panic(panicPrefix + err.Error())
	}
}
