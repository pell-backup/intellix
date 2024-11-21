package main

import (
	"flag"
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "world\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	var port string
	flag.StringVar(&port, "port", "8090", "port to serve on")
	flag.Parse()

	fmt.Printf("Serving on :%s...\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
