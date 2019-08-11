package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/majiru/ramble"
	"github.com/majiru/ramble/internal/uuid"
)

func handleSendHelloReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}

	var req ramble.SendHelloReq

	if json.Unmarshal(b, &req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}

	// TODO (esote): This should be done in the core.
	if req.Message == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return
	}

	// TODO (esote): This should be done in the core.
	if len(req.Recipients) == 0 {
		http.Error(w, "Empty recipient list", http.StatusBadRequest)
		return
	}

	// TODO (esote): This should be done in the core. (Where the core also
	// needs to check if Conversation doesn't exist, in which case give an
	// error.
	if req.Conversation == "" {
		if req.Conversation, err = uuid.UUID(); err != nil {
			http.Error(w, "Couldn't generate UUID", http.StatusInternalServerError)
			return
		}
	}

	// TODO (majiru): call core.

	var resp ramble.SendHelloResp

	n, err := ramble.NewHelloResponse()

	if err != nil {
		http.Error(w, "Couldn't generate resp", http.StatusInternalServerError)
		return
	}

	resp = ramble.SendHelloResp(*n)

	if b, err = json.Marshal(&resp); err != nil {
		http.Error(w, "Couldn't marshal send hello response", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}

// TODO (esote): handleSendVerifyReq, where core is called

func handleViewHelloReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}

	var req ramble.ViewHelloReq

	if json.Unmarshal(b, &req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}

	// TODO (esote): call core

	var resp ramble.ViewHelloResp

	n, err := ramble.NewHelloResponse()

	if err != nil {
		http.Error(w, "Couldn't generate resp", http.StatusInternalServerError)
		return
	}

	resp = ramble.ViewHelloResp(*n)

	if b, err = json.Marshal(&resp); err != nil {
		http.Error(w, "Couldn't marshal view hello response", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}

func handleViewVerifyReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Method", http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}

	var req ramble.ViewVerifyReq

	if json.Unmarshal(b, &req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}

	// TODO (esote): from core: look up UUID in mapping of active UUIDs &
	// verify time within 1 minute.
	if req.UUID == "" {
		http.Error(w, "Bad UUID", http.StatusBadRequest)
		return
	}

	//TODO(majiru): call core.

	var resp = ramble.ViewVerifyResp{Messages: nil}

	if b, err = json.Marshal(&resp); err != nil {
		http.Error(w, "Couldn't marshal view verify response", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}

func handleView(w http.ResponseWriter, r *http.Request) {
	// TODO (esote): fix how path is gathered: r.URL.String()[1:] string
	// until last substring with slash, etc.
	switch path.Base(r.URL.String()) {
	case "hello":
		handleViewHelloReq(w, r)
	case "verify":
		handleViewVerifyReq(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
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

	// TODO(majiru): Call core, nuke messages based on initial request type
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	// TODO (esote): fix how path is gathered: r.URL.String()[1:] string
	// until last substring with slash, etc.
	switch path.Base(r.URL.String()) {
	case "handle":
		handleDeleteHelloReq(w, r)
	case "verify":
		handleDeleteVerifyReq(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}
