package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	rootFs := http.FileServer(http.Dir("./"))

	artifactsHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		rootFs.ServeHTTP(resp, req)
		log.Println("artifactsHandler: ", req.URL.Path)
	})

	http.Handle("/", artifactsHandler)

	fmt.Println("Starting http server... Example url: http://0.0.0.0:5050/. Check the console!")
	if err := http.ListenAndServe(":5050", nil); err != nil {
		fmt.Println(err)
	}
}
