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
	wasmPath := filepath.Join(examplePath, "../wasm")
	artifactsPath := filepath.Join(examplePath, "../artifacts")

	exampleFs := http.FileServer(http.Dir(examplePath))
	wasmFs := http.FileServer(http.Dir(wasmPath))
	artifactsFs := http.FileServer(http.Dir(artifactsPath))

	http.Handle("/", exampleFs)
	http.Handle("/wasm/", http.StripPrefix("/wasm/", wasmFs))
	http.Handle("/artifacts/", http.StripPrefix("/artifacts/", artifactsFs))

	fmt.Println("Starting http server... Example url: http://localhost:8080/. Check the console!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
