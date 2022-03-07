package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"regexp"
	"time"
)

func getFileDate(obj *s3.HeadObjectOutput, filename string, filepattern string) time.Time {
	/*
		it returns the date of the file or the date within the filename if filepattern exists and matches
	*/
	filePattern, err := regexp.Compile(filepattern)
	if err == nil {
		r := filePattern.FindStringSubmatch(filename)
		if len(r) > 1 {
			theDate, err := time.Parse("2006-01-02", r[1])
			if err == nil {
				return theDate
			}
		}
	}
	return *obj.LastModified
}
