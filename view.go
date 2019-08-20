package ramble

const (
	// ViewConversations asks to view a list of conversations you are
	// associated with.
	ViewConversations uint8 = iota

	// ViewMessages asks to view a list messages within a conversion.
	ViewMessages
)

// ViewHelloReq is sent by the client as the initial request to view a list of
// stored messages.
type ViewHelloReq struct {
	// Count of how many items to return, 0 for all.
	Count uint64 `json:"count"`

	// Sender's public key fingerprint.
	Sender string `json:"sender"`

	// Type of data to view, representing an enumerated type.
	Type uint8 `json:"type"`
}

// ViewHelloResp is sent by the server in response to ViewHelloReq.
type ViewHelloResp HelloResponse

// ViewVerifyReq is sent by the client in response to ViewHelloResp.
type ViewVerifyReq VerifyRequest

// ViewVerifyResp is sent by the server in response to ViewVerifyReq and
// terminates the hello-verify handshake. The list items are encrypted with the
// sender's public key in an amalgamated string using speculative key IDs.
type ViewVerifyResp struct {
	// List of data.
	List string `json:"list"`
}
