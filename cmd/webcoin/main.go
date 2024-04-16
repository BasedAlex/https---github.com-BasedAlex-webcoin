package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const addr = "80"

func main() {
	srv := &http.Server{
		Addr:              ":" + addr,
		Handler:           routes(),
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Println(srv.Addr)
	err := srv.ListenAndServe()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", ping)

	return mux
}

func ping(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}
