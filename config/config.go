package config

import (
	"io/ioutil"
	"path"
	"strings"
	"os"
	"os/user"
	"github.com/jessevdk/go-flags"
	"fmt"

)

type Options struct {
    AccessKey string `short:"a" long:"access-key" description:"AccessKey"`
    SecretKey string `short:"s" long:"secret-key" description:"SecretKey"`
    Bucket    string `short:"b" long:"bucket" description:"s3 bucket"`
    Prefix    string `long:"prefix" description:"prefix of s3 location"`
    Output    string `short:"o" long:"output" description:"output file name"`
}

func Parse(opts *Options, args []string)([]string, error){
	fmt.Println(opts)
	fmt.Println(args)
	positionArgs, err := flags.ParseArgs(&opts, args)

	return positionArgs, err
}

func setup_opts(opts *Options) bool {
	if opts.AccessKey == "" {
		opts.AccessKey = FileConfig("access_key")
	}
	if opts.AccessKey == "" {
		opts.AccessKey = os.Getenv("S3CL_ACCESS_KEY")
	}
	if opts.AccessKey == "" {
		fmt.Println("access key not found")
		return false
	}


	if opts.SecretKey == "" {
		opts.SecretKey = FileConfig("secret_key")
	}
	if opts.SecretKey == "" {
		opts.SecretKey = os.Getenv("S3CL_SECRET_KEY")
	}
	if opts.SecretKey == "" {
		fmt.Println("secret key not found")
		return false
	}

	if opts.Bucket == "" {
		opts.Bucket = FileConfig("bucket")
	}
	if opts.Bucket == "" {
		opts.Bucket = os.Getenv("S3CL_BUCKET_DEFAULT")
	}
	if opts.Bucket == "" {
		fmt.Println("bucket not found")
		return false
	}

	return true
}

func FileConfig(key string)(val string) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	config_file := path.Join(usr.HomeDir, ".s3cl")
	content, err := ioutil.ReadFile(config_file)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
	    panic(err)
	}
	
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		vals := strings.Split(line, "=")
		if len(vals) != 2 {
			continue
		}
		if strings.TrimSpace(vals[0]) == key {
			return vals[1]
		}
	}

	return ""
}