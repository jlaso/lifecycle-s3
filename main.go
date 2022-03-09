package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
		cl := s3.New(sess, &aws.Config{
			Region:                         &cfg.Region,
			DisableRestProtocolURICleaning: aws.Bool(true),
		})
		req, resp := cl.ListObjectsRequest(&s3.ListObjectsInput{
			Bucket:    &cfg.Bucket,
			Delimiter: aws.String("/"),
			Prefix:    &cfg.Prefix,
		})

		var toKeep bool
		err = req.Send()
		if err == nil { // resp is now filled
			for i, s := range resp.Contents {
				toKeep = false
				if verbose {
					fmt.Println("************************************************************")
					//fmt.Printf("%d: %s %s\n", i, s.LastModified.Format("2006-01-02"), *s.Key)
				}
				fmt.Printf("%d: %s %s\n", i, s.LastModified.Format("2006-01-02"), *s.Key)
				for n, r := range cfg.Rules {
					fi := fileInfo{Name: *s.Key, Date: *s.LastModified}
					fi.Date = getFileDate(fi, cfg.FilePattern)
					toKeep = keepIt(r.Rule, fi)
					if verbose {
						fmt.Printf("\t%t\t%s\t%d\t%s %s\n",
							toKeep,
							(fi.Date).Weekday().String()[0:3],
							fileAge(fi.Date),
							n,
							r.Rule,
						)
					}
					if toKeep {
						break
					}
				}
				if toKeep {
					fmt.Println("\t\t   keep it !")
				} else {
					if cfg.Mode == "move-to-trash" {
						fmt.Println("\t\t   moving it to _TRASH_ !")
						if !sandbox {
							err = moveToTrash(cl, cfg.Bucket, *s.Key)
						}
					} else if cfg.Mode == "mark-with-tag" {
						fmt.Println("\t\t   marking it for deletion !")
						if !sandbox {
							err = markForDeletion(cl, cfg.Bucket, *s.Key)
						}
					}
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
	fmt.Fprintf(os.Stderr, "usage: lifecycle-s3 configfiles...\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.BoolVar(&verbose, "verbose", true, "show intermediate info")
	flag.BoolVar(&sandbox, "sandbox", true, "show intermediate info")
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
