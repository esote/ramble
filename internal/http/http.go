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
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}
	req := &ramble.SendReq{}
	if json.Unmarshal(b, req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	if req.Sender == "" {
		http.Error(w, "Empty sender", http.StatusBadRequest)
		return
	}
	//TODO(majiru): We also need to check that the msg is armor formated pgp message
	if req.Msg == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return
	}
	if req.Recipient == nil || len(req.Recipient) < 1 {
		http.Error(w, "Empty recipient list", http.StatusBadRequest)
		return
	}
	if req.Guid == 0 {
		//Generate new GUID
	}

	//TODO(majiru): Store msg

	//At this point req.Guid is filled with the correct Guid
	b, err = json.Marshal(&ramble.SendResp{req.Guid})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func handleView1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}
	req := &ramble.ViewReq1{}
	if json.Unmarshal(b, req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	if req.Fingerprint == "" {
		http.Error(w, "Empty fingerprint", http.StatusBadRequest)
		return
	}
	resp := &ramble.ViewResp1{}
	//TODO(majiru): populate resp and note time
	b, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func handleView2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Method", http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}
	req := &ramble.ViewReq2{}
	if json.Unmarshal(b, req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	//TODO(majiru): Verify time
	if req.SignedNONCE == "" {
		http.Error(w, "Empty NONCE", http.StatusBadRequest)
		return
	}
	if req.Guid == 0 {
		http.Error(w, "Invalid GUID", http.StatusBadRequest)
		return
	}
	//TODO(majiru): Grab convos
	convos := []int{}
	resp := &ramble.ViewResp2{convos}
	b, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func handleView(w http.ResponseWriter, r *http.Request) {
	switch path.Base(r.URL.String()) {
	case "1":
		handleView1(w, r)
	case "2":
		handleView2(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func handleDelete1(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}
	req := &ramble.DeleteReq1{}
	if json.Unmarshal(b, req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	if req.Fingerprint == "" {
		http.Error(w, "Empty fingerprint", http.StatusBadRequest)
		return
	}
	resp := &ramble.DeleteResp{}
	//TODO(majiru): populate resp and note time
	b, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func handleDelete2(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusBadRequest)
		return
	}
	req := &ramble.DeleteReq2{}
	if json.Unmarshal(b, req) != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	//TODO(majiru): Verify time
	if req.SignedNONCE == "" {
		http.Error(w, "Empty NONCE", http.StatusBadRequest)
		return
	}
	if req.Guid == 0 {
		http.Error(w, "Invalid GUID", http.StatusBadRequest)
		return
	}
	//TODO(majiru): Nuke messages
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleDelete1(w, r)
	case http.MethodGet:
		handleDelete2(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}
