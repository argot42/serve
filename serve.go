package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <host>:<port> <path>\n", os.Args[0])
		return
	}

	http.Handle("/", http.FileServer(http.Dir(os.Args[2])))
	log.Fatal(http.ListenAndServe(os.Args[1], nil))
}
