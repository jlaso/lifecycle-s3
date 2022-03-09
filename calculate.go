package main

import (
	"fmt"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	"log"
	"reflect"
	"time"
)

type fileInfo struct {
	Name string
	Date time.Time
}

func keepIt(condition string, file fileInfo) bool {
	e := env.NewEnv()

	err := e.Define("file_age", fileAge(file.Date))
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}
	err = e.Define("file_time", file.Date)
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}
	err = e.Define("first_day_of_year", isFirstDayOfYear)
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}
	err = e.Define("first_day_of_month", isFirstDayOfMonth)
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}
	err = e.Define("last_day_of_year", isLastDayOfYear)
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}
	err = e.Define("last_day_of_month", isLastDayOfMonth)
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}
	err = e.Define("last_day_of_week", isLastDayOfWeek)
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}

	result, err := vm.Execute(e, nil, fmt.Sprintf("return (%s)", condition))
	if err != nil {
		log.Fatalf("Execute error `%s` (%s): %v\n", condition, file, err)
	}

	return reflect.ValueOf(result).Bool()
}
