package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path"
	"strings"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"config"
)



var opts = config.Options{"", "", "", "", ""}

func help() {
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

func main() {
	args, err := config.Parse(&opts, os.Args)
	if err != nil {
    	panic(err)
    	os.Exit(1)
	}

	if len(args) < 2 {
		help()	
		return
	}

	//fmt.Println(opts)
	subcmd := args[1]
	if subcmd == "help" {
		help()
	} else if subcmd == "config" {
		CmdConf()
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



func CmdConf() {
	fmt.Println("todo")
}

func CmdGet(opts config.Options, key string) {
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

func Get(opts config.Options, key string) (data []byte, err error){
	auth := aws.Auth{AccessKey: opts.AccessKey, SecretKey: opts.SecretKey}
	s3Instance := s3.New(auth, aws.Region{Name: "*", S3Endpoint: "http://s3.amazonaws.com"})
	bkt := s3Instance.Bucket(opts.Bucket)

	return bkt.Get(key)
}


func CmdPut(opts config.Options, key string, file_paths []string) {

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

func Put(opts config.Options, key string, data []byte, contType string) error {
	auth := aws.Auth{AccessKey: opts.AccessKey, SecretKey: opts.SecretKey}
	s3Instance := s3.New(auth, aws.Region{Name: "*", S3Endpoint: "http://s3.amazonaws.com"})	
	bkt := s3Instance.Bucket(opts.Bucket)

	return bkt.Put(key, data, contType, s3.BucketOwnerFull, s3.Options{})
}

func CmdMove(opts config.Options, src_key, dest_key string) {
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

func Del(opts config.Options, key string) error{
	auth := aws.Auth{AccessKey: opts.AccessKey, SecretKey: opts.SecretKey}
	s3Instance := s3.New(auth, aws.Region{Name: "*", S3Endpoint: "http://s3.amazonaws.com"})
	bkt := s3Instance.Bucket(opts.Bucket)

	return bkt.Del(key)
}

func CmdDel(opts config.Options, key string) {
	err := Del(opts, key)
	if err != nil {
		panic(err)
	}
}

func List(opts config.Options) error{
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

func CmdList(opts config.Options) {
	err := List(opts)
	if err != nil {
		panic(err)
	}
}