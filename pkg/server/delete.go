package server

import (
	"errors"
	"strings"

	"github.com/majiru/ramble"
	"github.com/majiru/ramble/internal/pgp"
)

// DeleteHello processes the hello handshake step.
func (s *Server) DeleteHello(req *ramble.DeleteHelloReq) (*ramble.DeleteHelloResp, error) {
	if !pgp.VerifyHexFingerprint(req.Sender) {
		return nil, errors.New("sender fingerprint is invalid")
	}

	req.Sender = strings.ToLower(req.Sender)

	resp, err := s.newHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.DeleteHelloResp(*resp)

	return &ret, nil
}

// DeleteVerify processes the verify handshake step.
func (s *Server) DeleteVerify(req *ramble.DeleteVerifyReq) (*ramble.DeleteVerifyResp, error) {
	meta, err := s.verifyReq(req.UUID)

	if err != nil {
		return nil, err
	}

	hello, ok := meta.request.(*ramble.DeleteHelloReq)

	if !ok {
		return nil, errors.New("request was not DeleteHelloReq")
	}

	public, err := s.public.Read(hello.Sender)

	if err != nil {
		return nil, err
	}

	if err = s.verifyReqSig(public, req.Signature, meta.nonce); err != nil {
		return nil, err
	}

	// TODO: delete things based on hello.Type

	return new(ramble.DeleteVerifyResp), nil
}
