package ramble

import (
	"github.com/majiru/ramble/internal/pgp"
	"github.com/majiru/ramble/internal/uuid"
)

func NewHelloResponse() (*HelloResponse, error) {
	var h HelloResponse

	b, err := pgp.NonceHex()

	if err != nil {
		return nil, err
	}

	h.Nonce = string(b)

	u, err := uuid.UUID()

	if err != nil {
		return nil, err
	}

	h.UUID = u

	return &h, nil
}
