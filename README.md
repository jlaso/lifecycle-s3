# lifecycle-s3

This software is aim to solve a lack in s3 lifecycle, that is the fact of removing file according to a
plan based on the date of the files.

In order to do so, a configuration file in toml format is needed, that's the reason of using
- https://github.com/BurntSushi/toml

This configuration file should be passed as a parameter when invoking the executable, in fact several of this
config files can be passed separated by spaces.

This config file looks like that:

    bucket = "your-bucket-name"  <-  the bucket name
    region = "us-east-1"    <- here the region where the bucket is
    prefix = "folder-to-review/subfolder/etc/"    <- the folder to process
    mode = "move-to-trash"  <- can be mark-to-delete
    filepattern = ""   <-  fill this in if you prefer that the date is deduced from the file name instead
    
    [rules]
    ...             <- here the rules to apply in order to keep or discard the file


## additional parameters

As usual, you can invoke the command with -verbose in order to see the intermediate steps

Also, a -sandbox mode can be enabled to not mark or delete in S3.


## rules description

The parameters or functions ready to use in a rule are:
    
    - file_time  <- the date of the file (real or deduced from file name)
    - file_age  <-  the days that the file is old (based on file_time)
    - first_day_of_year(d)  <- true if date passed is Jan 1st
    - first_day_of_month(d)   <- true if date passed is 1st of month
    - last_day_of_year(d)   <-  true if date passed is Dec 31st
    - last_day_of_month(d)   <- true if date passed is the end of month
    - last_day_of_week(d)  <-  true if date passed is Saturday

The functions and variables available can be upgraded easily (check calculate.go)

### Example

 
    [rules.keep-last-of-year-2] <- this is just a name, should start by rules. 
    # <--- this rule will keep files that were created by Dec 31st --->
    rule = "last_day_of_year(file_time)" 
    
    [rules.keep-monthly-12]
    # <--- this rule will keep end of month files within the current 365 period --->
    rule = "(file_age <= 365) && last_day_of_month(file_time)"
    
    [rules.keep-8-weeks]
    # <--- this rule will keep Saturday files within the current 8 weeks period --->
    rule = "(file_age <= 56) && last_day_of_week(file_time)"
    
    [rules.keep-last-7-days]
    # <--- this rule will keep last 7 days files --->
    rule = "(file_age <= 7)"


Since the idea is executing this software once a day (for instance), the rules should 
consider the rotation of files. When a file is aged out, if an old-rule should keep it, 
the youngest one should consider this keeping as well.

### modes

So far the current modes supported are:

- move-to-trash
- mark-with-tag
- delete (unsupported)

    #### * move-to-trash
    
    The discardable file is moved to a _TRASH_ folder withing the processing folder.
    After the process a human can review this folder and empty it or just apply S3 lifecycles 
    over it.
    
    #### * mark-with-tag
    
    A DELETE_ME tag is added to the discardable. So far this option seems to be unuseful since
    S3 cannot delete files based on tags.
    
    ##### * delete
    
    Currently, not supported since the software is in beta stage.
    
    
