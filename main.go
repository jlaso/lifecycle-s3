package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func checkFiles(filename string) error {

	cfg, err := readConfig(filename)

	fmt.Printf("bucket name is %s, found %d rules \n", cfg.Bucket, len(cfg.Rules))
	for i, a := range cfg.Rules {
		fmt.Printf("rule %s: %s \n", i, a.Rule)
	}

	if len(cfg.Rules) > 0 {

		sess := session.Must(session.NewSession())
		//client := s3manager.
		cl := s3.New(sess, &aws.Config{
			Region:                         &cfg.Region,
			DisableRestProtocolURICleaning: aws.Bool(true),
		})
		delimiter := "/"
		req, resp := cl.ListObjectsRequest(&s3.ListObjectsInput{
			Bucket:    &cfg.Bucket,
			Delimiter: &delimiter,
			Prefix:    &cfg.Prefix,
		})

		var toKeep bool
		err = req.Send()
		if err == nil { // resp is now filled
			//fmt.Println(resp.Contents)
			for i, s := range resp.Contents {
				toKeep = false
				//fmt.Println("************************************************************")
				//fmt.Printf("%d: %s %s\n", i, s.LastModified.Format("2006-01-02"), *s.Key)
				fmt.Printf("%d: %s %s", i, s.LastModified.Format("2006-01-02"), *s.Key)
				//for n, r := range cfg.Rules {
				for _, r := range cfg.Rules {
					toKeep = toKeep || check(r.Rule, fileInfo{Name: *s.Key, Date: *s.LastModified})
					//fmt.Printf("\t%t\t%d\t%d\t%s %s\n",
					//	!toKeep,
					//	(*s.LastModified).Weekday(),
					//	fileAge(*s.LastModified),
					//	n,
					//	r.Rule,
					//)
				}
				if toKeep {
					fmt.Println("   keep it !")
				} else {
					fmt.Println("   marking it for deletion !")
					err = markForDeletion(cl, cfg.Bucket, *s.Key)
					if err != nil {
						return err
					}
				}
			}
		} else {
			return err
		}
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: lifecycle-s3 configfile\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}
	for _, f := range args {
		err := checkFiles(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}
