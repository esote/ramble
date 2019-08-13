package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/majiru/ramble"
)

func handleView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	switch path.Base(r.URL.String()) {
	case "hello":
		handleViewHelloReq(w, r)
	case "verify":
		handleViewVerifyReq(w, r)
	default:
		writeError(w, http.StatusNotFound)
	}
}

func handleViewHelloReq(w http.ResponseWriter, r *http.Request) {
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

	resp, err := ramble.ViewHello(&req)

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

func handleViewVerifyReq(w http.ResponseWriter, r *http.Request) {
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

	resp, err := ramble.ViewVerify(&req)

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
