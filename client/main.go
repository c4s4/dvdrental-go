package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	URL = "http://localhost:8080/"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("ERROR: pass url and count (such as actor and 200) on command line")
		os.Exit(-1)
	}
	url := os.Args[1]
	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	for i := 1; i <= count; i++ {
		response, err := http.Get(URL + url + "/" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		println(string(content))
	}
}
