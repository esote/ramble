package ramble

import (
	"github.com/majiru/ramble/internal/pgp"
	"github.com/majiru/ramble/internal/uuid"
)

// HelloResponse is sent from the server indicating that it needs verification
// before continuing.
type HelloResponse struct {
	// Nonce to be signed and passed to the verify request.
	Nonce string `json:"nonce"`

	// UUID to be passed to the verify request.
	UUID string `json:"uuid"`
}

// NewHelloResponse generates a hello response with nonce and UUID.
func NewHelloResponse() (*HelloResponse, error) {
	var h HelloResponse

	b, err := pgp.NonceHex()

	if err != nil {
		return nil, err
	}

	h.Nonce = string(b)

	h.UUID, err = uuid.UUID()

	if err != nil {
		return nil, err
	}

	return &h, nil
}

// VerifyRequest is sent from the client with verification details. The
// signature is used to verify ownership of a private key.
type VerifyRequest struct {
	// Signature is the detached signature of the hello response nonce.
	Signature string `json:"sig"`

	// UUID from the hello response.
	UUID string `json:"uuid"`
}
