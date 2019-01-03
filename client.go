package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	URL = "http://localhost:8080/actor/"
	Max = 200
)

func main() {
	for i := 1; i < Max; i++ {
		response, err := http.Get(URL + strconv.Itoa(i))
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
