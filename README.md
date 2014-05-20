# Command line tools for Amazon S3 services.


## Target
* fast
* out of box
* easy to use
* shell friendly

## Expected usage
```
$s3cl help
$s3cl config

$s3cl get key
$s3cl put key value(file, stdin, string, number)
$s3cl mv  key_src key_dest(file?/directory?)
$s3cl del key
```

## Configuration search order
1. command line options, like ```--access-key --secret-key```
2. config file located in ```~/.s3cl```
3. environment variables like ```S3CL_ACCESS_KEY, S3CL_SECRET_KEY```