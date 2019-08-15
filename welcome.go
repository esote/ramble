package ramble

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/majiru/ramble/internal/pgp"
)

// WelcomeHelloReq is sent from the client asking to add this public key to
// storage. This is required before all other requests, since all other requests
// initiate based on the sender's fingerprint, not full public key.
type WelcomeHelloReq struct {
	// Public key.
	Public string `json:"public"`
}

// WelcomeHelloResp is sent by the server in response to WelcomeHelloReq.
type WelcomeHelloResp HelloResponse

// WelcomeVerifyReq is sent by the client in response to WelcomeHelloResp.
type WelcomeVerifyReq VerifyRequest

// WelcomeVerifyResp is sent by  the server in response to WelcomeVerifyReq and
// terminates the hello-verify handshake.
type WelcomeVerifyResp struct{}

// WelcomeHello processes the hello handshake step.
func WelcomeHello(w *WelcomeHelloReq) (*WelcomeHelloResp, error) {
	public := strings.NewReader(w.Public)

	if ok, err := pgp.VerifyPublicArmored(public); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("input not a public key")
	}

	resp, err := NewHelloResponse(w)

	if err != nil {
		return nil, err
	}

	ret := WelcomeHelloResp(*resp)

	return &ret, nil
}

// WelcomeVerify processes the verify handshake step.
func WelcomeVerify(w *WelcomeVerifyReq) (*WelcomeVerifyResp, error) {
	m, ok := activeHVs[w.UUID]

	if !ok {
		return nil, errors.New("no handshake with UUID")
	}

	delete(activeHVs, w.UUID)

	if time.Now().UTC().Sub(m.time) > maxHVDur {
		return nil, errors.New("handshake expired")
	}

	hello, ok := m.request.(*WelcomeHelloReq)

	if !ok {
		return nil, errors.New("request was not WelcomeHelloReq")
	}

	public := strings.NewReader(hello.Public)
	sig := strings.NewReader(w.Signature)
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

	err = pub.Write(hex.EncodeToString(fingerprint), []byte(hello.Public))

	if err != nil {
		return nil, err
	}

	return new(WelcomeVerifyResp), nil
}
