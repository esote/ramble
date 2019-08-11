package ramble

import (
	"time"

	"github.com/majiru/ramble/internal/pgp"
	"github.com/majiru/ramble/internal/uuid"
)

func init() {
	// Launch globally-persisting goroutine used to prune handshakes older
	// than MaxHVDur. The activeHVs time value should still be checked since
	// this cannot remove stale handshakes immediately.
	go func() {
		ticker := time.NewTicker(maxHVDur)

		for {
			select {
			case now := <-ticker.C:
				for uuid, m := range activeHVs {
					if now.UTC().Sub(m.time) > maxHVDur {
						delete(activeHVs, uuid)
					}
				}
			}
		}
	}()
}

// Max allowed time between the hello request and the expected verify response.
const maxHVDur = time.Minute

// Metadata used when verifying and proceeding in a hello-verify handshake.
type verifyMeta struct {
	nonce   string
	request interface{}
	time    time.Time
}

// Active hello-verify handshakes.
var activeHVs = make(map[string]verifyMeta)

// HelloResponse is sent from the server indicating that it needs verification
// before continuing.
type HelloResponse struct {
	// Nonce to be signed and passed to the verify request.
	Nonce string `json:"nonce"`

	// UUID to be passed to the verify request.
	UUID string `json:"uuid"`
}

// NewHelloResponse generates a hello response with nonce and UUID, then adds it
// to the active hello-verify handshakes map.
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
