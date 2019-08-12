package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/majiru/ramble"
)

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	switch path.Base(r.URL.String()) {
	case "hello":
		handleSendHelloReq(w, r)
	case "verify":
		handleSendVerifyReq(w, r)
	default:
		writeError(w, http.StatusNotFound)
	}
}

func handleSendHelloReq(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.SendHelloReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := ramble.SendHello(&req)

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

func handleSendVerifyReq(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	var req ramble.SendVerifyReq

	if json.Unmarshal(b, &req) != nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	resp, err := ramble.SendVerify(&req)

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
