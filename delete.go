package ramble

import (
	"errors"
	"strings"
	"time"

	"github.com/majiru/ramble/internal/pgp"
)

// DeleteHelloReq is sent by the client as the initial request to delete stored
// data.
type DeleteHelloReq struct {
	// Sender's public key fingerprint.
	Sender string `json:"sender"`

	// Type of data to delete, representing an enumerated type.
	Type uint8 `json:"type"`
}

// DeleteHelloResp is sent by the server in response to DeleteHelloReq.
type DeleteHelloResp HelloResponse

// DeleteVerifyReq is sent by the client in response to DeleteHelloResp.
type DeleteVerifyReq VerifyRequest

// DeleteVerifyResp is sent by the server in response to DeleteVerifyReq and
// terminates the hello-verify handshake.
type DeleteVerifyResp struct{}

// TODO: these have no meaning yet
const (
	DeleteAll uint8 = iota
	DeleteReceived
	DeleteSent

	deleteMax = DeleteSent
)

// DeleteHello processes the hello handshake step.
func DeleteHello(d *DeleteHelloReq) (*DeleteHelloResp, error) {
	if d.Type > deleteMax {
		return nil, errors.New("invalid delete type")
	}

	resp, err := NewHelloResponse(d)

	if err != nil {
		return nil, err
	}

	ret := DeleteHelloResp(*resp)

	return &ret, nil
}

// DeleteVerify processes the verify handshake step.
func DeleteVerify(d *DeleteVerifyReq) (*DeleteVerifyResp, error) {
	m, ok := activeHVs[d.UUID]

	if !ok {
		return nil, errors.New("no handshake with UUID")
	}

	delete(activeHVs, d.UUID)

	if time.Now().UTC().Sub(m.time) > maxHVDur {
		return nil, errors.New("handshake expired")
	}

	hello, ok := m.request.(*DeleteHelloReq)

	if !ok {
		return nil, errors.New("request was not DeleteHelloReq")
	}

	// TODO: use hello.Sender to index for public key
	_ = hello
	public := strings.NewReader("TODO")

	sig := strings.NewReader(d.Signature)
	nonce := strings.NewReader(m.nonce)

	if ok, err := pgp.VerifyArmoredSig(public, sig, nonce); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("signature did not match public key")
	}

	// TODO: delete things based on hello.Type

	return new(DeleteVerifyResp), nil
}
