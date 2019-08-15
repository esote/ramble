package server

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/majiru/ramble"
	"github.com/majiru/ramble/internal/pgp"
)

// WelcomeHello processes the hello handshake step.
func (s *Server) WelcomeHello(req *ramble.WelcomeHelloReq) (*ramble.WelcomeHelloResp, error) {
	public := strings.NewReader(req.Public)

	if ok, err := pgp.VerifyPublicArmored(public); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("input not a public key")
	}

	resp, err := s.NewHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.WelcomeHelloResp(*resp)

	return &ret, nil
}

// WelcomeVerify processes the verify handshake step.
func (s *Server) WelcomeVerify(req *ramble.WelcomeVerifyReq) (*ramble.WelcomeVerifyResp, error) {
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

	hello, ok := m.request.(*ramble.WelcomeHelloReq)

	if !ok {
		return nil, errors.New("request was not WelcomeHelloReq")
	}

	public := strings.NewReader(hello.Public)
	sig := strings.NewReader(req.Signature)
	nonce := strings.NewReader(m.nonce)

	if ok, err := pgp.VerifyArmoredSig(public, sig, nonce); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("signature did not match public key")
	}

	// Reset reader.
	public = strings.NewReader(hello.Public)

	fingerprint, err := pgp.FingerprintArmored(public)

	if err != nil {
		return nil, errors.New("unable to get public key fingerprint")
	}

	err = s.public.Write(hex.EncodeToString(fingerprint), []byte(hello.Public))

	if err != nil {
		return nil, err
	}

	return new(ramble.WelcomeVerifyResp), nil
}
