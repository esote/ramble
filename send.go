package ramble

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/majiru/ramble/internal/pgp"
)

// SendHelloReq is sent by the client as the initial hello request to append a
// message to a conversion.
type SendHelloReq struct {
	// Conversation UUID representing a pre-existing conversation, or empty
	// to start a new conversation.
	Conversation string `json:"conv"`

	// Message PGP encrypted. The list of encryption recipients should match
	// the "recipients" member.
	Message string `json:"msg"`

	// Recipients' public key fingerprints.
	Recipients []string `json:"recipients"`

	// Sender's public key fingerprint.
	Sender string `json:"sender"`
}

// SendHelloResp is sent by the server in response to SendHelloReq.
type SendHelloResp HelloResponse

// SendVerifyReq is sent by the client in response to SendHelloResp.
type SendVerifyReq VerifyRequest

// SendVerifyResp is sent by the server in response to SendVerifyReq and
// terminates the hello-verify handshake.
type SendVerifyResp struct{}

// SendHello processes the hello handshake step.
func SendHello(s *SendHelloReq) (*SendHelloResp, error) {
	if len(s.Recipients) == 0 {
		return nil, errors.New("empty recipient list")
	}

	msg := strings.NewReader(s.Message)

	if ok, err := pgp.VerifyEncryptedArmored(msg); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("message is not encrypted and armored")
	}

	response, err := NewHelloResponse()

	if err != nil {
		return nil, err
	}

	if m, ok := activeHVs[response.UUID]; ok {
		log.Printf("send: %s -> %s already exists in activeHVs!\n",
			response.UUID, m.time.String())
		return nil, errors.New("the very improbable just happened")
	}

	activeHVs[response.UUID] = verifyMeta{
		nonce:   response.Nonce,
		request: s,
		time:    time.Now().UTC(),
	}

	ret := SendHelloResp(*response)

	return &ret, nil
}

// SendVerify processes the verify handshake step.
func SendVerify(s *SendVerifyReq) (*SendVerifyResp, error) {
	m, ok := activeHVs[s.UUID]

	if !ok {
		return nil, errors.New("no handshake with UUID")
	}

	delete(activeHVs, s.UUID)

	if time.Now().UTC().Sub(m.time) > maxHVDur {
		return nil, errors.New("handshake expired")
	}

	hello, ok := m.request.(*SendHelloReq)

	if !ok {
		return nil, errors.New("request was not SendHelloReq")
	}

	// TODO: use hello.Sender to index for public key
	_ = hello
	public := strings.NewReader("TODO")

	sig := strings.NewReader(s.Signature)
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

	return new(SendVerifyResp), nil
}
