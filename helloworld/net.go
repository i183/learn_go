package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello, World!")
	})
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello, World!2")
	})
	go startWeb(":1234")
	go startWeb(":1235")
	time.Sleep(1000000000000000)
}

func startWeb(addr string) {

	http.ListenAndServe(addr, nil)
}
