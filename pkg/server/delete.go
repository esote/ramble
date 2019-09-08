package server

import (
	"errors"
	"strings"

	"github.com/esote/ramble"
	"github.com/esote/ramble/internal/pgp"
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

	switch hello.Type {
	case ramble.DeleteAll:
		err = s.public.Remove(hello.Sender)

		if err == nil {
			err = s.tconvos.Splay.Remove(hello.Sender)
		}
	case ramble.DeletePublic:
		err = s.public.Remove(hello.Sender)
	case ramble.DeleteConversations:
		err = s.tconvos.Splay.Remove(hello.Sender)
	}

	if err != nil {
		return nil, err
	}

	return new(ramble.DeleteVerifyResp), nil
}
