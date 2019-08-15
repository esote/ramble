package server

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/majiru/ramble"
	"github.com/majiru/ramble/internal/pgp"
)

// DeleteHello processes the hello handshake step.
func (s *Server) DeleteHello(req *ramble.DeleteHelloReq) (*ramble.DeleteHelloResp, error) {
	if !pgp.VerifyHexFingerprint(req.Sender) {
		return nil, errors.New("sender fingerprint is invalid")
	}

	req.Sender = strings.ToLower(req.Sender)

	resp, err := s.NewHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.DeleteHelloResp(*resp)

	return &ret, nil
}

// DeleteVerify processes the verify handshake step.
func (s *Server) DeleteVerify(req *ramble.DeleteVerifyReq) (*ramble.DeleteVerifyResp, error) {
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

	hello, ok := m.request.(*ramble.DeleteHelloReq)

	if !ok {
		return nil, errors.New("request was not DeleteHelloReq")
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

	// TODO: delete things based on hello.Type

	return new(ramble.DeleteVerifyResp), nil
}
