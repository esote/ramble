package server

import (
	"errors"
	"strings"

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

	resp, err := s.newHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.ViewHelloResp(*resp)

	return &ret, nil
}

// ViewVerify processes the verify handshake step.
func (s *Server) ViewVerify(req *ramble.ViewVerifyReq) (*ramble.ViewVerifyResp, error) {
	meta, err := s.verifyReq(req.UUID)

	if err != nil {
		return nil, err
	}

	hello, ok := meta.request.(*ramble.ViewHelloReq)

	if !ok {
		return nil, errors.New("request was not ViewHelloReq")
	}

	public, err := s.public.Read(hello.Sender)

	if err != nil {
		return nil, err
	}

	if err = s.verifyReqSig(public, req.Signature, meta.nonce); err != nil {
		return nil, err
	}

	// TODO: return index of messages sent by user and sent to user

	// TODO: this return definitely needs to be encrypted using their public
	// key, so no one else can read the data.

	return new(ramble.ViewVerifyResp), nil
}
