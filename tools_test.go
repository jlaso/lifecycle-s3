package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDateInFileName(t *testing.T) {
	dateToTest := time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)
	filePattern := "source-(\\d{4}-\\d{2}-\\d{2})\\.zip"
	obj := fileInfo{
		Name: fmt.Sprintf("path/subpath/chetl-source-%s.zip", dateToTest.Format("2006-01-02")),
		Date: time.Now(),
	}
	dt := getFileDate(obj, filePattern)
	if dt != dateToTest {
		t.Errorf("the returned date is not the expected one %s", dateToTest.Format("2006-01-02"))
	}
}

func TestDateInFileObj(t *testing.T) {
	dateToTest := time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)
	filePattern := ""
	today := time.Now()
	obj := fileInfo{
		Name: fmt.Sprintf("path/subpath/chetl-source-%s.zip", dateToTest.Format("2006-01-02")),
		Date: today,
	}
	dt := getFileDate(obj, filePattern)
	if dt != today {
		t.Errorf("the returned date is not the expected one %s", today.Format("2006-01-02"))
	}
}
