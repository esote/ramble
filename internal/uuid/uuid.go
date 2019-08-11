package uuid

import (
	"encoding/hex"

	"github.com/esote/util/uuid"
)

// UUID generates a new UUID string.
func UUID() (string, error) {
	u, err := uuid.NewUUID()

	if err != nil {
		return "", err
	}

	h := make([]byte, hex.EncodedLen(len(u)))
	_ = hex.Encode(h, u)

	return string(h), nil
}
