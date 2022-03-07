package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"path"
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

func moveToTrash(c *s3.S3, bucket string, fname string) error {
	srcKey := "/" + bucket + "/" + fname
	destKey := path.Dir(fname) + "/_TRASH_/" + path.Base(fname)
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fname),
	}
	sourceObject, err := c.HeadObject(headInput)
	if err != nil {
		return fmt.Errorf("failed to read object head: %v", err)
	}
	fmt.Println(sourceObject)
	meta := sourceObject.Metadata
	if meta == nil {
		meta = make(map[string]*string)
	}
	meta["modified"] = aws.String(sourceObject.LastModified.Format("2006-01-02"))
	_, err = c.CopyObject(
		&s3.CopyObjectInput{
			Bucket:            aws.String(bucket),
			Key:               aws.String(destKey),
			CopySource:        aws.String(srcKey),
			MetadataDirective: aws.String("REPLACE"),
			Metadata:          meta,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to copy object: %v", err)
	}
	_, err = c.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fname),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}
	return nil
}
