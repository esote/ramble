package ramble

import "time"

// TODO (esote): add a welcome struct JSON also using hello-verify handshake
// which requests we store their public key. Using this, we no longer need them
// to pass the full public key in send, verify, and delete, but instead just the
// public key fingerprint.

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

// HelloResponse is sent from the server indicating that it needs verification
// before continuing.
type HelloResponse struct {
	// Nonce to be signed and passed to the verify request.
	Nonce string `json:"nonce"`

	// UUID to be passed to the verify request.
	UUID string `json:"uuid"`
}

// VerifyRequest is sent from the client with verification details. The
// signature is used to verify ownership of a private key.
type VerifyRequest struct {
	// Signature is the detached signature of the hello response nonce.
	Signature string `json:"sig"`

	// UUID from the hello response.
	UUID string `json:"uuid"`
}

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

	// Sender's public key.
	Sender string `json:"sender"`
}

// SendHelloResp is sent by the server in response to SendHelloReq.
type SendHelloResp HelloResponse

// SendVerifyReq is sent by the client in response to SendHelloResp.
type SendVerifyReq VerifyRequest

// SendVerifyResp is sent by the server in response to SendVerifyReq and
// terminates the hello-verify handshake.
type SendVerifyResp struct{}

// ViewHelloReq is sent by the client as the initial request to view a list of
// stored messages.
type ViewHelloReq struct {
	// Sender's public key.
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

// DeleteHelloReq is sent by the client as the initial request to delete stored
// data.
type DeleteHelloReq struct {
	// Sender's public key.
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
