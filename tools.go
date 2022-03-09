package main

import (
	"regexp"
	"time"
)

func getFileDate(fi fileInfo, filepattern string) time.Time {
	/*
		it returns the date of the file or the date within the filename if filepattern exists and matches
	*/
	filePattern, err := regexp.Compile(filepattern)
	if err == nil {
		r := filePattern.FindStringSubmatch(fi.Name)
		if len(r) > 1 {
			theDate, err := time.Parse("2006-01-02", r[1])
			if err == nil {
				return theDate
			}
		}
	}
	return fi.Date
}
