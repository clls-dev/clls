package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	go func() {
		if err := changeWatcher(os.Getenv("PWD")); err != nil {
			fmt.Printf("watcher error: %s\n", err)
		}
	}()
	if err := http.ListenAndServe(`127.0.0.1:8080`, http.FileServer(http.Dir(`.`))); err != nil {
		fmt.Printf("serve error: %s\n", err)
	}
}
