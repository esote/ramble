package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/majiru/ramble/pkg/server"
)

var srv *server.Server

func main() {
	var err error
	srv, err = server.NewServer(time.Hour, "test_convo", "test_msg", "test_public")

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	path := filepath.Clean(r.URL.Path)

	switch path {
	case "/delete/hello":
		handleDeleteHello(w, r)
	case "/delete/verify":
		handleDeleteVerify(w, r)
	case "/send/hello":
		handleSendHello(w, r)
	case "/send/verify":
		handleSendVerify(w, r)
	case "/view/hello":
		handleViewHello(w, r)
	case "/view/verify":
		handleViewVerify(w, r)
	case "/welcome/hello":
		handleWelcomeHello(w, r)
	case "/welcome/verify":
		handleWelcomeVerify(w, r)
	default:
		writeError(w, http.StatusNotFound)
	}
}
