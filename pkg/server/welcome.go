package server

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/esote/ramble"
	"github.com/esote/ramble/internal/pgp"
)

// WelcomeHello processes the hello handshake step.
func (s *Server) WelcomeHello(req *ramble.WelcomeHelloReq) (*ramble.WelcomeHelloResp, error) {
	public := strings.NewReader(req.Public)

	if ok, err := pgp.VerifyPublicArmored(public); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("input not a public key")
	}

	resp, err := s.newHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.WelcomeHelloResp(*resp)

	return &ret, nil
}

// WelcomeVerify processes the verify handshake step.
func (s *Server) WelcomeVerify(req *ramble.WelcomeVerifyReq) (*ramble.WelcomeVerifyResp, error) {
	meta, err := s.verifyReq(req.UUID)

	if err != nil {
		return nil, err
	}

	hello, ok := meta.request.(*ramble.WelcomeHelloReq)

	if !ok {
		return nil, errors.New("request was not WelcomeHelloReq")
	}

	err = s.verifyReqSig([]byte(hello.Public), req.Signature, meta.nonce)

	if err != nil {
		return nil, err
	}

	public := strings.NewReader(hello.Public)
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
