package main

import (
	"net/http"
	"strings"
)

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.String()[1:], "/")

	if len(path) < 2 {
		writeError(w, http.StatusNotFound)
		return
	}

	// Only care about end of path, ignore any leading prefix.
	path = path[len(path)-2:]

	switch path[0] {
	case "delete":
		switch path[1] {
		case "hello":
			handleDeleteHello(w, r)
		case "verify":
			handleDeleteVerify(w, r)
		default:
			writeError(w, http.StatusNotFound)
		}
	case "send":
		switch path[1] {
		case "hello":
			handleSendHello(w, r)
		case "verify":
			handleSendVerify(w, r)
		default:
			writeError(w, http.StatusNotFound)
		}
	case "view":
		switch path[1] {
		case "hello":
			handleViewHello(w, r)
		case "verify":
			handleViewVerify(w, r)
		default:
			writeError(w, http.StatusNotFound)
		}
	case "welcome":
		switch path[1] {
		case "hello":
			handleWelcomeHello(w, r)
		case "verify":
			handleWelcomeVerify(w, r)
		default:
			writeError(w, http.StatusNotFound)
		}
	default:
		writeError(w, http.StatusNotFound)
	}
}
