package uuid

import (
	"encoding/hex"

	"github.com/esote/util/uuid"
)

// LenUUID is the hex length of a UUID.
const LenUUID = 32

// UUID generates a new UUID string.
func UUID() (string, error) {
	u, err := uuid.NewUUID()

	if err != nil {
		return "", err
	}

	h := make([]byte, LenUUID)
	_ = hex.Encode(h, u)

	return string(h), nil
}
