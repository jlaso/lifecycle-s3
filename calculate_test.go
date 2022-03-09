package main

import (
	"testing"
	"time"
)

func newFileInfo(m time.Month, d int) fileInfo {
	today := time.Now()
	return fileInfo{
		Name: "test",
		Date: time.Date(today.Year(), m, d, 0, 0, 0, 0, time.UTC),
	}
}

func TestCheck(t *testing.T) {

	var battery = []struct {
		fi       fileInfo
		cond     string
		expected bool
	}{
		{newFileInfo(1, 1), "first_day_of_year(file_time)", true},
		{newFileInfo(1, 2), "first_day_of_year(file_time)", false},
		{newFileInfo(12, 31), "last_day_of_year(file_time)", true},
		{newFileInfo(3, 7), "last_day_of_year(file_time)", false},
		{newFileInfo(3, 7), "(file_age > 365) && last_day_of_year(file_time)", false},
		{newFileInfo(3, 7), "(file_age > 365)", false},
	}

	for i, b := range battery {
		if keepIt(b.cond, b.fi) != b.expected {
			t.Errorf("%d: Unexpected result for '%s'", i, b.cond)
		}
	}

}
