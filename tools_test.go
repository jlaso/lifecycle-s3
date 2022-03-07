package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"testing"
	"time"
)

func TestDateInFileName(t *testing.T) {
	dateToTest := time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)
	key := fmt.Sprintf("path/subpath/chetl-source-%s.zip", dateToTest.Format("2006-01-02"))
	filePattern := "source-(\\d{4}-\\d{2}-\\d{2})\\.zip"
	obj := s3.HeadObjectOutput{
		LastModified: aws.Time(time.Now()),
	}
	dt := getFileDate(&obj, key, filePattern)
	if dt != dateToTest {
		t.Errorf("the returned date is not the expected one %s", dateToTest.Format("2006-01-02"))
	}
}

func TestDateInFileObj(t *testing.T) {
	dateToTest := time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)
	key := fmt.Sprintf("path/subpath/chetl-source-%s.zip", dateToTest.Format("2006-01-02"))
	filePattern := ""
	today := time.Now()
	obj := s3.HeadObjectOutput{
		LastModified: aws.Time(today),
	}
	dt := getFileDate(&obj, key, filePattern)
	if dt != today {
		t.Errorf("the returned date is not the expected one %s", today.Format("2006-01-02"))
	}
}
