package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"path"
)

type S3Instance struct {
	S3Client *s3.S3
}

func (s *S3Instance) connect(region string) {
	sess := session.Must(session.NewSession())
	s.S3Client = s3.New(sess, &aws.Config{
		Region:                         aws.String(region),
		DisableRestProtocolURICleaning: aws.Bool(true),
	})
}

func (s *S3Instance) walkFiles(bucket string, prefix string, c chan *fileInfo) error {
	var contToken *string
	var err error
	for {
		resp, err := s.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			ContinuationToken: contToken,
			Delimiter:         aws.String("/"),
			Prefix:            aws.String(prefix),
		})
		if err == nil {
			for _, s := range resp.Contents {
				c <- &fileInfo{Name: *s.Key, Date: *s.LastModified}
			}
		}
		if resp.ContinuationToken == nil {
			break
		}
		contToken = resp.ContinuationToken
	}
	c <- nil
	return err
}

func (s *S3Instance) markForDeletion(bucket string, fname string) error {

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
	_, err := s.S3Client.PutObjectTagging(&poti)

	return err
}

func (s *S3Instance) moveToTrash(bucket string, fname string) error {
	srcKey := "/" + bucket + "/" + fname
	destKey := path.Dir(fname) + "/_TRASH_/" + path.Base(fname)
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fname),
	}
	sourceObject, err := s.S3Client.HeadObject(headInput)
	if err != nil {
		return fmt.Errorf("failed to read object head: %v", err)
	}
	fmt.Println(sourceObject)
	meta := sourceObject.Metadata
	if meta == nil {
		meta = make(map[string]*string)
	}
	meta["modified"] = aws.String(sourceObject.LastModified.Format("2006-01-02"))
	_, err = s.S3Client.CopyObject(
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
	_, err = s.S3Client.DeleteObject(
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
