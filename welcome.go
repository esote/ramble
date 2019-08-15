package ramble

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
