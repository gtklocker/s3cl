package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)


func main() {
	fmt.Println("curl setup")
	resp, err := http.Get("http://jatsz.org/")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}