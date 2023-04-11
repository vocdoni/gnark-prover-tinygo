package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("unable to get the current caller file")
		return
	}
	examplePath := filepath.Dir(callerFile)
	artifactsPath := filepath.Join(examplePath, "../../../artifacts")

	exampleFs := http.FileServer(http.Dir(examplePath))
	artifactsFs := http.FileServer(http.Dir(artifactsPath))

	artifactsHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		artifactsFs.ServeHTTP(resp, req)
		log.Println("artifactsHandler: ", req.URL.Path)
	})

	http.Handle("/", exampleFs)
	http.Handle("/artifacts/", http.StripPrefix("/artifacts/", artifactsHandler))

	fmt.Println("Starting http server... Example url: http://0.0.0.0:5050/. Check the console!")
	if err := http.ListenAndServe(":5050", nil); err != nil {
		fmt.Println(err)
	}
}
