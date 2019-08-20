package server

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/esote/ramble"
	"github.com/esote/ramble/internal/pgp"
	"github.com/esote/ramble/internal/uuid"
)

var reHex = regexp.MustCompile("^[a-fA-F0-9]+$")

// SendHello processes the hello handshake step.
func (s *Server) SendHello(req *ramble.SendHelloReq) (*ramble.SendHelloResp, error) {
	if len(req.Recipients) == 0 {
		return nil, errors.New("empty recipient list")
	}

	if req.Conversation == "" {
		var err error
		req.Conversation, err = uuid.UUID()

		if err != nil {
			return nil, err
		}
	}

	if len(req.Conversation) != uuid.LenUUID ||
		!reHex.MatchString(req.Conversation) {
		return nil, errors.New("conversation UUID invalid")
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

		req.Recipients[i] = strings.ToLower(r)
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

	msg, err := uuid.UUID()

	if err != nil {
		return nil, err
	}

	if err = s.tmsgs.Insert(hello.Conversation, msg); err != nil {
		return nil, err
	}

	if err = s.msg.Write(msg, []byte(hello.Message)); err != nil {
		return nil, err
	}

	err = s.tconvos.InsertUnique(hello.Sender, hello.Conversation)

	if err != nil {
		return nil, err
	}

	for _, r := range hello.Recipients {
		err = s.tconvos.InsertUnique(r, hello.Conversation)

		if err != nil {
			return nil, err
		}
	}

	return &ramble.SendVerifyResp{
		Conversation: hello.Conversation,
	}, nil
}
