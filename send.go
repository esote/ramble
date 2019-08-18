package ramble

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
type SendVerifyResp struct {
	// Conversation UUID. If the hello request conversation UUID was empty,
	// this UUID is for the new conversation.
	Conversation string `json:"conv"`
}
