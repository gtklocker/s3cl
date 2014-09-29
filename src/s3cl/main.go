package main

import (
	"fmt"
	"os"
	"os/user"
	"io/ioutil"
	"path"
	"strings"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/jessevdk/go-flags"
)

type Options struct {
    AccessKey string `short:"a" long:"access-key" description:"AccessKey"`
    SecretKey string `short:"s" long:"secret-key" description:"SecretKey"`
    Bucket    string `short:"b" long:"bucket" description:"s3 bucket"`
    Prefix    string `long:"prefix" description:"prefix of s3 location"`
    Output    string `short:"o" long:"output" description:"output file name"`
}

var opts = Options{}


func main() {
	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
    	panic(err)
    	os.Exit(1)
	}

	if len(args) < 2 {
		help()	
		return
	}

	if !setup_opts(&opts) {
		return
	}

	//fmt.Println(opts)

	subcmd := args[1]
	if subcmd == "help" {
		help()
	} else if subcmd == "config" {
		config()
	} else if subcmd == "get" {
		CmdGet(opts, args[2])
	} else if subcmd == "put" {
		CmdPut(opts, args[2], args[3:])
	} else if subcmd == "del" {
		CmdDel(opts, args[2])
	} else if subcmd == "mv" {
		CmdMove(opts, args[2], args[3])
	} else if subcmd == "ls" {
		CmdList(opts)
	} else {  //default
		help()
	}

	return
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

func print_usage() {
	usage := `
Usage:
    s3cl get key
    s3cl put key file_path
    s3cl del key
    s3cl mv  src_key dest_key
    
Options:
    --access-key, -a    Access Key
    --secret-key, -s    Secret Key
    --bucket    , -b    Bucket Name
    --prefix            Prefix of s3 location
    --output    , -o    Output file name
    `
	fmt.Println(usage)
}

func help() {
	print_usage()
}

func config() {
	
}

func CmdGet(opts Options, key string) {
	d, err := Get(opts, key)
	if err != nil {
		panic(err)
	}
	if opts.Output != "" {
		err := ioutil.WriteFile(opts.Output, d, 0666)
		if err != nil {
	    	panic(err)
		}
	} else {
		os.Stdout.Write(d)
		os.Stdout.Sync()
	}
}

func Get(opts Options, key string) (data []byte, err error){
	auth := aws.Auth{AccessKey: opts.AccessKey, SecretKey: opts.SecretKey}
	s3Instance := s3.New(auth, aws.Region{Name: "*", S3Endpoint: "http://s3.amazonaws.com"})
	bkt := s3Instance.Bucket(opts.Bucket)

	return bkt.Get(key)
}


func CmdPut(opts Options, key string, file_paths []string) {

	if strings.HasSuffix(key, "/") {
		for _, file_path := range file_paths {

			contType := "text/plain"
			if strings.HasSuffix(file_path, ".xml") {
				contType = "text/xml"
			} else if strings.HasSuffix(file_path,".flv") {
				contType = "video/x-flv"
			}

			d, err := ioutil.ReadFile(file_path)
			if err != nil {
				panic(err)
				return
			}

			err = Put(opts, key + path.Base(file_path), d, contType)
			if err != nil {
				panic(err)
				return
			}
		}
	} else {
		//if not ends with /, just put the first file
		file_path := file_paths[0]
		contType := "text/plain"
		if strings.HasSuffix(file_path, ".xml") {
			contType = "text/xml"
		} else if strings.HasSuffix(file_path,".flv") {
			contType = "video/x-flv"
		}

		d, err := ioutil.ReadFile(file_path)
		if err != nil {
			panic(err)
			return
		}

		err = Put(opts, key, d, contType)
		if err != nil {
			panic(err)
			return
		}
	}


}

func Put(opts Options, key string, data []byte, contType string) error {
	auth := aws.Auth{AccessKey: opts.AccessKey, SecretKey: opts.SecretKey}
	s3Instance := s3.New(auth, aws.Region{Name: "*", S3Endpoint: "http://s3.amazonaws.com"})	
	bkt := s3Instance.Bucket(opts.Bucket)

	return bkt.Put(key, data, contType, s3.BucketOwnerFull, s3.Options{})
}

func CmdMove(opts Options, src_key, dest_key string) {
	d, err := Get(opts, src_key)
	if err != nil {
		panic(err)
	}

	Put(opts, dest_key, d, "")
	if err != nil {
		panic(err)
		return
	}
	Del(opts, src_key)
}

func Del(opts Options, key string) error{
	auth := aws.Auth{AccessKey: opts.AccessKey, SecretKey: opts.SecretKey}
	s3Instance := s3.New(auth, aws.Region{Name: "*", S3Endpoint: "http://s3.amazonaws.com"})
	bkt := s3Instance.Bucket(opts.Bucket)

	return bkt.Del(key)
}

func CmdDel(opts Options, key string) {
	err := Del(opts, key)
	if err != nil {
		panic(err)
	}
}

func List(opts Options) error{
	auth := aws.Auth{AccessKey: opts.AccessKey, SecretKey: opts.SecretKey}
	s3Instance := s3.New(auth, aws.Region{Name: "*", S3Endpoint: "http://s3.amazonaws.com"})
	bkt := s3Instance.Bucket(opts.Bucket)

	lst, err := bkt.List(opts.Prefix, "", "", 100)
	// fmt.Println(lst)
	if err == nil {
		for _, r := range lst.Contents {
			fmt.Println(r.Key)
		}
	}

	return err
}

func CmdList(opts Options) {
	err := List(opts)
	if err != nil {
		panic(err)
	}
}