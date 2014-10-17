# Command line tools for Amazon S3 services.

## Status
[ ![Codeship Status for jatsz/s3cl](https://codeship.io/projects/8b60ee50-37c9-0132-96d5-02255209da1c/status)](https://codeship.io/projects/41810)

## Target
* fast
* out of box
* easy to use
* shell friendly

## Expected usage
```
$s3cl help
$s3cl config

$s3cl ls  s3://{bucket_name}{key}
$s3cl get s3://{bucket_name}{key}
$s3cl put s3://{bucket_name}{key} local_file/stdin
$s3cl mv  s3://{bucket_name}{key} s3://{bucket_name}{key}
$s3cl del s3://{bucket_name}{key}
```

## Configuration search order
1. command line options, like ```--access-key --secret-key```
2. config file located in ```~/.s3cl```
3. environment variables like ```S3CL_ACCESS_KEY, S3CL_SECRET_KEY```