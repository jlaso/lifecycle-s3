package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

func markForDeletion(c *s3.S3, bucket string, fname string) error {

	deleteMe := "DELETE_ME"
	tag := s3.Tagging{
		TagSet: []*s3.Tag{
			&s3.Tag{
				Key:   &deleteMe,
				Value: &deleteMe,
			},
		},
	}
	poti := s3.PutObjectTaggingInput{
		Bucket:  &bucket,
		Key:     &fname,
		Tagging: &tag,
	}
	_, err := c.PutObjectTagging(&poti)

	return err
}
