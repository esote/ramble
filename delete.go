package ramble

// TODO: these have no meaning yet
const (
	DeleteAll uint8 = iota
	DeletePublic
	DeleteConversations
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
