package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/majiru/ramble"
)

func handleViewHello(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.ViewHelloReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := srv.ViewHello(&req)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	if b, err = json.Marshal(resp); err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}

func handleViewVerify(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.ViewVerifyReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := srv.ViewVerify(&req)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	if b, err = json.Marshal(resp); err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}
