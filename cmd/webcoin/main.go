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
	fmt.Println(srv.Addr)
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

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	fmt.Fprintln(w, "Hello world!")
}