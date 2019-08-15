package main

import (
	"log"
	"net/http"
	"time"

	"github.com/majiru/ramble/pkg/server"
)

var srv *server.Server

func main() {
	var err error
	srv, err = server.NewServer(time.Hour, "test_publickeys")

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
