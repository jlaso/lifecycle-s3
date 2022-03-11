package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func checkFiles(filename string) error {

	cfg, err := readConfig(filename)
	if err != nil {
		return err
	}

	fmt.Printf("bucket name is %s, found %d rules \n", cfg.Bucket, len(cfg.Rules))
	for i, a := range cfg.Rules {
		fmt.Printf("rule %s: %s \n", i, a.Rule)
	}

	if len(cfg.Rules) > 0 {
		var myS3 S3Instance
		myS3.connect(cfg.Region)

		ch := make(chan *fileInfo)

		go func() {
			var keepIt bool
			i := 1
			for {
				fi := <-ch
				if fi == nil {
					break
				}
				keepIt = false
				if verbose {
					fmt.Println("************************************************************")
				}
				fmt.Printf("%d: %s %s\n", i, fi.Date.Format("2006-01-02"), fi.Name)
				for n, r := range cfg.Rules {
					fi.Date = getFileDate(*fi, cfg.FilePattern)
					keepIt = fi.canBeKept(r.Rule)
					if verbose {
						fmt.Printf("\t%t\t%s\t%d\t%s %s\n",
							keepIt,
							(fi.Date).Weekday().String()[0:3],
							fileAge(fi.Date),
							n,
							r.Rule,
						)
					}
					if keepIt {
						break
					}
				}
				if keepIt {
					fmt.Println("\t\t   keep it !")
				} else {
					if cfg.Mode == "move-to-trash" {
						fmt.Println("\t\t   moving it to _TRASH_ !")
						if !sandbox {
							err = myS3.moveToTrash(cfg.Bucket, fi.Name)
						}
					} else if cfg.Mode == "mark-with-tag" {
						fmt.Println("\t\t   marking it for deletion !")
						if !sandbox {
							err = myS3.markForDeletion(cfg.Bucket, fi.Name)
						}
					}
				}
				i++
			}
		}()

		err = myS3.walkFiles(cfg.Bucket, cfg.Prefix, ch)

	}
	return err
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
		fmt.Println("Config file is missing.")
		os.Exit(1)
	}
	for _, f := range args {
		err := checkFiles(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}
