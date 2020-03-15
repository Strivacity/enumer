// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Marshaling test: test marshaler intefaces.

package main

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

type SQL int

const (
	Monday SQL = iota
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
	ckUnmarshal("Monday", Monday)
	ckUnmarshal("Tuesday", Tuesday)
	ckUnmarshal("Wednesday", Wednesday)
	ckUnmarshal("Thursday", Thursday)
	ckUnmarshal("Friday", Friday)
	ckUnmarshal("Saturday", Saturday)
	ckUnmarshal("Sunday", Sunday)
	ckUnmarshal("Christmas", 127)
}

const panicPrefix = "sql.go: "

func ckNoError(err error) {
	if err != nil {
		panic(panicPrefix + err.Error())
	}
}

func ckMarshal(e SQL) {
	db, mock, err := sqlmock.New()
	ckNoError(err)
	defer db.Close()
	mock.ExpectExec("INSERT INTO marshal_test").WithArgs(e.String()).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = db.Exec("INSERT INTO marshal_test (sql_enum) VALUES (?,)", e)
	ckNoError(err)
	err = mock.ExpectationsWereMet()
	ckNoError(err)
}

func ckUnmarshal(raw string, expected SQL) {
	db, mock, err := sqlmock.New()
	ckNoError(err)
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"sql_enum"}).AddRow(raw))
	result, err := db.Query("SELECT")
	ckNoError(err)
	defer result.Close()
	if !result.Next() {
		panic(panicPrefix + "expected row")
	}
	got := SQL(-1)
	err = result.Scan(&got)
	if !expected.IsASQL() {
		if err == nil {
			panic(panicPrefix + "expected error")
		}
	} else {
		ckNoError(err)
		if got != expected {
			panic(fmt.Sprintf("%s sql.Scan got '%s', expected '%s'", panicPrefix, got, expected))
		}
	}
	if result.Next() {
		panic(panicPrefix + "expected no more rows")
	}
}
