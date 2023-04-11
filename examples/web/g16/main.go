package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
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

	http.Handle("/", exampleFs)
	http.Handle("/artifacts/", http.StripPrefix("/artifacts/", artifactsFs))

	fmt.Println("Starting http server... Example url: http://0.0.0.0:5050/. Check the console!")
	if err := http.ListenAndServe(":5050", nil); err != nil {
		fmt.Println(err)
	}
}
