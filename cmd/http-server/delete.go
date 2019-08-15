package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/majiru/ramble"
)

func handleDeleteHello(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.DeleteHelloReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := ramble.DeleteHello(&req)

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

func handleDeleteVerify(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.DeleteVerifyReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := ramble.DeleteVerify(&req)

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
