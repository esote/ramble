package server

import (
	"bytes"
	"errors"
	"strings"

	"github.com/esote/ramble"
	"github.com/esote/ramble/internal/pgp"
	"github.com/esote/ramble/internal/uuid"
)

// ViewHello processes the hello handshake step.
func (s *Server) ViewHello(req *ramble.ViewHelloReq) (*ramble.ViewHelloResp, error) {
	switch req.Type {
	case ramble.ViewConversations, ramble.ViewMessages:
		break
	default:
		return nil, errors.New("invalid type")
	}

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

	var buf bytes.Buffer

	switch hello.Type {
	case ramble.ViewConversations:
		convos, err := s.tconvos.IndexN(hello.Sender, hello.Count)

		if err != nil {
			return nil, err
		}

		buf.Grow(len(convos) * (uuid.LenUUID + 1))

		for _, convo := range convos {
			buf.WriteString(convo)
			buf.Write([]byte{'\n'})
		}
	case ramble.ViewMessages:
		msgs, err := s.tmsgs.IndexN(hello.Sender, hello.Count)

		if err != nil {
			return nil, err
		}

		for _, msgUUID := range msgs {
			msg, err := s.msg.Read(msgUUID)

			if err != nil {
				return nil, err
			}

			buf.Write(msg)
			buf.Write([]byte{'\n'})
		}
	default:
		return nil, errors.New("invalid type")
	}

	p := bytes.NewReader(public)
	l := bytes.NewReader(buf.Bytes())

	enc, err := pgp.EncryptArmored(p, l)

	if err != nil {
		return nil, err
	}

	return &ramble.ViewVerifyResp{
		List: string(enc),
	}, nil
}
