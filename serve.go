package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	nargs := len(os.Args)
	if nargs < 2 {
		usage()
	}

	var verbose bool
	var path string
	var socket string

	path = "."

	switch os.Args[1] {
	case "-h":
		usage()
	case "-v":
		verbose = true
		socket = os.Args[2]
		if nargs > 3 {
			path = os.Args[3]
		}
	default:
		socket = os.Args[1]
		if nargs > 2 {
			path = os.Args[2]
		}
	}

	handler := http.FileServer(http.Dir(path))
	if verbose {
		handler = logger(handler)
	}
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(socket, nil))
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-hv] [<host>]:<port> [<path>]\n", os.Args[0])
	os.Exit(1)
}
