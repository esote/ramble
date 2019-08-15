package server

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/majiru/ramble"
	"github.com/majiru/ramble/internal/pgp"
)

// SendHello processes the hello handshake step.
func (s *Server) SendHello(req *ramble.SendHelloReq) (*ramble.SendHelloResp, error) {
	if len(req.Recipients) == 0 {
		return nil, errors.New("empty recipient list")
	}

	if !pgp.VerifyHexFingerprint(req.Sender) {
		return nil, errors.New("sender fingerprint is invalid")
	}

	req.Sender = strings.ToLower(req.Sender)

	for i, r := range req.Recipients {
		if !pgp.VerifyHexFingerprint(r) {
			return nil, fmt.Errorf("recipient fingerprint index=%d"+
				" is invalid", i)
		}

		r = strings.ToLower(r)
	}

	msg := strings.NewReader(req.Message)

	if ok, err := pgp.VerifyEncryptedArmored(msg); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("message is not encrypted and armored")
	}

	resp, err := s.NewHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.SendHelloResp(*resp)

	return &ret, nil
}

// SendVerify processes the verify handshake step.
func (s *Server) SendVerify(req *ramble.SendVerifyReq) (*ramble.SendVerifyResp, error) {
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

	hello, ok := m.request.(*ramble.SendHelloReq)

	if !ok {
		return nil, errors.New("request was not SendHelloReq")
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

	// TODO: store hello.Message as part of (or creating) m.Conversation
	// with the metadata hello.Recipients.

	// TODO: maybe force that messages can only be sent to recipients who
	// have added their own public keys through the welcome process?

	return new(ramble.SendVerifyResp), nil
}
