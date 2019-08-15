package ramble

import (
	"time"
)

// StoredMessage is a message and its metadata.
type StoredMessage struct {
	// Conversation UUID.
	Conversation string `json:"conv"`

	// Encrypted message.
	Message string `json:"msg"`

	// Recipients' public key fingerprints.
	Recipients []string `json:"recipients"`

	// Time the message was sent.
	Time time.Time `json:"time"`
}

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
