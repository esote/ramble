package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/majiru/ramble"
)

func writeStatus(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeStatus(http.StatusMethodNotAllowed, w)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	req := &ramble.SendReq{}
	if json.Unmarshal(b, req) != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	if req.Sender == "" {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	//TODO(majiru): We also need to check that the msg is armor formated pgp message
	if req.Msg == "" {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	if req.Recipient == nil || len(req.Recipient) < 1 {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	if req.Guid == 0 {
		//Generate new GUID
	}

	//TODO(majiru): Store msg

	//At this point req.Guid is filled with the correct Guid
	b, err = json.Marshal(&ramble.SendResp{req.Guid})
	if err != nil {
		writeStatus(http.StatusInternalServerError, w)
		return
	}
	w.Write(b)
}

func handleView1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeStatus(http.StatusMethodNotAllowed, w)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	req := &ramble.ViewReq1{}
	if json.Unmarshal(b, req) != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	if req.Fingerprint == "" {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	resp := &ramble.ViewResp1{}
	//TODO(majiru): populate resp and note time
	b, err = json.Marshal(resp)
	if err != nil {
		writeStatus(http.StatusInternalServerError, w)
		return
	}
	w.Write(b)
}

func handleView2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeStatus(http.StatusMethodNotAllowed, w)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	req := &ramble.ViewReq2{}
	if json.Unmarshal(b, req) != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	//TODO(majiru): Verify time
	if req.SignedNONCE == "" {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	if req.Guid == 0 {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	//TODO(majiru): Grab convos
	convos := []int{}
	resp := &ramble.ViewResp2{convos}
	b, err = json.Marshal(resp)
	if err != nil {
		writeStatus(http.StatusInternalServerError, w)
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
		writeStatus(http.StatusNotFound, w)
	}
}

func handleDelete1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeStatus(http.StatusMethodNotAllowed, w)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	req := &ramble.DeleteReq1{}
	if json.Unmarshal(b, req) != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	if req.Fingerprint == "" {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	resp := &ramble.DeleteResp{}
	//TODO(majiru): populate resp and note time
	b, err = json.Marshal(resp)
	if err != nil {
		writeStatus(http.StatusInternalServerError, w)
		return
	}
	w.Write(b)
}

func handleDelete2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeStatus(http.StatusMethodNotAllowed, w)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	req := &ramble.DeleteReq2{}
	if json.Unmarshal(b, req) != nil {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	//TODO(majiru): Verify time
	if req.SignedNONCE == "" {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	if req.Guid == 0 {
		writeStatus(http.StatusBadRequest, w)
		return
	}
	//TODO(majiru): Nuke messages
	writeStatus(http.StatusOK, w)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	switch path.Base(r.URL.String()) {
	case "1":
		handleDelete1(w, r)
	case "2":
		handleDelete2(w, r)
	default:
		writeStatus(http.StatusNotFound, w)
	}
}
