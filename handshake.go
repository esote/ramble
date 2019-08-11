package ramble

import (
	"errors"
	"log"
	"time"

	"github.com/majiru/ramble/internal/pgp"
	"github.com/majiru/ramble/internal/uuid"
)

// MaxHVDuration is the max allowed time between the hello request and the
// expected verify response.
const MaxHVDuration = time.Minute

func init() {
	activeHVs = make(map[string]time.Time)

	// Launch globally-persisting goroutine used to prune handshakes older
	// than MaxHVDuration. The activeHVs value should still be checked each
	// time since this cannot remove stale handshakes immediately.
	go func() {
		ticker := time.NewTicker(MaxHVDuration)

		for {
			select {
			case now := <-ticker.C:
				for uuid, t := range activeHVs {
					if now.Sub(t) > MaxHVDuration {
						delete(activeHVs, uuid)
					}
				}
			}
		}
	}()
}

// Map (UUID -> Time added) of active hello-verify handshakes.
var activeHVs map[string]time.Time

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

	if t, ok := activeHVs[h.UUID]; ok {
		log.Printf("UUID %s -> %s already exists in activeHVs!", h.UUID,
			t.String())
		return nil, errors.New("the very improbable just happened")
	}

	activeHVs[h.UUID] = time.Now()

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
