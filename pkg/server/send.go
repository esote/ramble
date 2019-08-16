package server

import (
	"errors"
	"fmt"
	"strings"

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

	resp, err := s.newHelloResponse(req)

	if err != nil {
		return nil, err
	}

	ret := ramble.SendHelloResp(*resp)

	return &ret, nil
}

// SendVerify processes the verify handshake step.
func (s *Server) SendVerify(req *ramble.SendVerifyReq) (*ramble.SendVerifyResp, error) {
	meta, err := s.verifyReq(req.UUID)

	if err != nil {
		return nil, err
	}

	hello, ok := meta.request.(*ramble.SendHelloReq)

	if !ok {
		return nil, errors.New("request was not SendHelloReq")
	}

	public, err := s.public.Read(hello.Sender)

	if err != nil {
		return nil, err
	}

	if err = s.verifyReqSig(public, req.Signature, meta.nonce); err != nil {
		return nil, err
	}

	// TODO: store hello.Message as part of (or creating) m.Conversation
	// with the metadata hello.Recipients.

	// TODO: maybe force that messages can only be sent to recipients who
	// have added their own public keys through the welcome process?

	return new(ramble.SendVerifyResp), nil
}
