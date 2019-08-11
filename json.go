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
