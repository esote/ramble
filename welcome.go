package ramble

import (
	"errors"
	"log"
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

	response, err := NewHelloResponse()

	if err != nil {
		return nil, err
	}

	if m, ok := activeHVs[response.UUID]; ok {
		log.Printf("welcome: %s -> %s already exists in activeHVs!\n",
			response.UUID, m.time.String())
		return nil, errors.New("the very improbable just happened")
	}

	activeHVs[response.UUID] = verifyMeta{
		nonce:   response.Nonce,
		request: w,
		time:    time.Now().UTC(),
	}

	ret := WelcomeHelloResp(*response)

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

	// TODO: store public key, if it already exists overwrite it.

	return new(WelcomeVerifyResp), nil
}
