package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"io"
	_ "embed"
)

//go:embed favicon.ico
var fav []byte

func main() {
	help := flag.Bool("h", false, "this help")
	verbose := flag.Bool("v", false, "verbose output")
	stdin := flag.Bool("i", false, "serve stdin")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	socket := ":1234"
	path := "."

	args := flag.Args()

	switch len(args) {
	case 1:
		socket = args[0]
	case 2:
		socket = args[0]
		path = args[1]
	}

	var handler http.Handler
	var initMsg string

	if *stdin {
		initMsg = "serving stdin to " + socket
		handler = servstr()
	} else {
		initMsg = "serving " + path + " to " + socket
		handler = http.FileServer(http.Dir(path))
	}
	if *verbose {
		handler = logger(handler)
		log.Println(initMsg)
	}

	http.HandleFunc("/favicon.ico", favHandler)
	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(socket, nil))
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func servstr() http.Handler {
	src, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(src)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func favHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(fav)
	if err != nil {
		log.Fatal(err)
	}
}
