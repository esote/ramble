package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/majiru/ramble"
)

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func handleDeleteHelloReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}

	var req ramble.DeleteHelloReq

	if json.Unmarshal(b, &req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}

	// TODO (esote): call core.

	var resp ramble.DeleteHelloResp

	n, err := ramble.NewHelloResponse()

	if err != nil {
		http.Error(w, "Couldn't generate resp", http.StatusInternalServerError)
		return
	}

	resp = ramble.DeleteHelloResp(*n)

	if b, err = json.Marshal(resp); err != nil {
		http.Error(w, "Couldn't marshal view verify response", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}

func handleDeleteVerifyReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}

	var req ramble.DeleteVerifyReq

	if json.Unmarshal(b, &req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}

	// TODO (majiru): Call core, nuke messages based on initial request type
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	switch path.Base(r.URL.String()) {
	case "hello":
		handleDeleteHelloReq(w, r)
	case "verify":
		handleDeleteVerifyReq(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}
