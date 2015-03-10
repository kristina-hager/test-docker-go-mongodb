package main

import (
	"fmt"
	"net/http"
	"runtime"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello there you beautiful world:  I was built on ", runtime.GOOS, " with a CPU ", runtime.GOARCH)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
