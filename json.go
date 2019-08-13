package ramble

import "time"

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

// DeleteHelloReq is sent by the client as the initial request to delete stored
// data.
type DeleteHelloReq struct {
	// Sender's public key fingerprint.
	Sender string `json:"sender"`

	// Type of data to delete, representing an enumerated type.
	Type uint64 `json:"type"`
}

// DeleteHelloResp is sent by the server in response to DeleteHelloReq.
type DeleteHelloResp HelloResponse

// DeleteVerifyReq is sent by the client in response to DeleteHelloResp.
type DeleteVerifyReq VerifyRequest

// DeleteVerifyResp is sent by the server in response to DeleteVerifyReq and
// terminates the hello-verify handshake.
type DeleteVerifyResp struct{}
