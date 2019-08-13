package ramble

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/majiru/ramble/internal/pgp"
)

// ViewHelloReq is sent by the client as the initial request to view a list of
// stored messages.
type ViewHelloReq struct {
	// Sender's public key fingerprint.
	Sender string `json:"sender"`

	// Count of how many messages to return.
	Count int64 `json:"count"`
}

// ViewHelloResp is sent by the server in response to ViewHelloReq.
type ViewHelloResp HelloResponse

// ViewVerifyReq is sent by the client in response to ViewHelloResp.
type ViewVerifyReq VerifyRequest

// ViewVerifyResp is sent by the server in response to ViewVerifyReq and
// terminates the hello-verify handshake.
type ViewVerifyResp struct {
	// Messages is a list of messages.
	Messages []StoredMessage `json:"msgs"`
}

// ViewHello processes the hello handshake step.
func ViewHello(v *ViewHelloReq) (*ViewHelloResp, error) {
	if v.Count <= 0 {
		return nil, errors.New("view count <= 0")
	}

	response, err := NewHelloResponse()

	if err != nil {
		return nil, err
	}

	if m, ok := activeHVs[response.UUID]; ok {
		log.Printf("view: %s -> %s already exists in activeHVs!\n",
			response.UUID, m.time.String())
		return nil, errors.New("the very improbable just happened")
	}

	activeHVs[response.UUID] = verifyMeta{
		nonce:   response.Nonce,
		request: v,
		time:    time.Now().UTC(),
	}

	ret := ViewHelloResp(*response)

	return &ret, nil
}

// ViewVerify processes the verify handshake step.
func ViewVerify(v *ViewVerifyReq) (*ViewVerifyResp, error) {
	m, ok := activeHVs[v.UUID]

	if !ok {
		return nil, errors.New("no handshake with UUID")
	}

	delete(activeHVs, v.UUID)

	if time.Now().UTC().Sub(m.time) > maxHVDur {
		return nil, errors.New("handshake expired")
	}

	hello, ok := m.request.(*ViewHelloReq)

	if !ok {
		return nil, errors.New("request was not ViewHelloReq")
	}

	// TODO: use hello.Sender to index for public key
	_ = hello
	public := strings.NewReader("TODO")

	sig := strings.NewReader(v.Signature)
	nonce := strings.NewReader(m.nonce)

	if ok, err := pgp.VerifyArmoredSig(public, sig, nonce); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("signature did not match public key")
	}

	// TODO: return index of messages sent by user and sent to user

	// TODO: this return definitely needs to be encrypted using their public
	// key, so no one else can read the data.

	return new(ViewVerifyResp), nil
}
