package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const (
	URL = "http://localhost:"
	Port = "8080"
	Size = 2
)

func feed(count int, index chan int) {
	for i := 1; i <= count; i++ {
		index <- i
	}
	close(index)
}

func work(operation string, index chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case i, ok := <-index:
			if ok {
				response, err := http.Get(URL + Port + "/" + operation + "/" + strconv.Itoa(i))
				if err != nil {
					panic(err)
				}
				content, err := ioutil.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}
				println(string(content))
			} else {
				println("Thread done")
				return
			}
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("ERROR: pass operation and count (such as actor and 200) on command line")
		os.Exit(-1)
	}
	operation := os.Args[1]
	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(Size)
	index := make(chan int)
	go feed(count, index)
	for i := 0; i < Size; i++ {
		go work(operation, index, &wg)
	}
	wg.Wait()
}
