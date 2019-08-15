package server

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/majiru/ramble"
	"github.com/majiru/ramble/internal/pgp"
)

// ViewHello processes the hello handshake step.
func (s *Server) ViewHello(req *ramble.ViewHelloReq) (*ramble.ViewHelloResp, error) {
	if !pgp.VerifyHexFingerprint(req.Sender) {
		return nil, errors.New("sender fingerprint is invalid")
	}

	req.Sender = strings.ToLower(req.Sender)

	if req.Count <= 0 {
		return nil, errors.New("view count <= 0")
	}

	resp, err := s.NewHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.ViewHelloResp(*resp)

	return &ret, nil
}

// ViewVerify processes the verify handshake step.
func (s *Server) ViewVerify(req *ramble.ViewVerifyReq) (*ramble.ViewVerifyResp, error) {
	s.mu.Lock()
	m, ok := s.active[req.UUID]

	if !ok {
		s.mu.Unlock()
		return nil, errors.New("no handshake with UUID")
	}

	delete(s.active, req.UUID)
	s.mu.Unlock()

	if time.Now().UTC().Sub(m.time) > s.dur {
		return nil, errors.New("handshake expired")
	}

	hello, ok := m.request.(*ramble.ViewHelloReq)

	if !ok {
		return nil, errors.New("request was not ViewHelloReq")
	}

	publicBody, err := s.public.Read(hello.Sender)

	if err != nil {
		return nil, err
	}

	public := bytes.NewReader(publicBody)
	sig := strings.NewReader(req.Signature)
	nonce := strings.NewReader(m.nonce)

	if ok, err := pgp.VerifyArmoredSig(public, sig, nonce); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("signature did not match public key")
	}

	// TODO: return index of messages sent by user and sent to user

	// TODO: this return definitely needs to be encrypted using their public
	// key, so no one else can read the data.

	return new(ramble.ViewVerifyResp), nil
}
