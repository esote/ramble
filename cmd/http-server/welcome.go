package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/majiru/ramble"
)

func handleWelcomeHello(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.WelcomeHelloReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := ramble.WelcomeHello(&req)

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

func handleWelcomeVerify(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.WelcomeVerifyReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := ramble.WelcomeVerify(&req)

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
