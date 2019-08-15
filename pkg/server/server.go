package server

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/esote/util/splay"
	"github.com/majiru/ramble"
	"github.com/majiru/ramble/internal/pgp"
	"github.com/majiru/ramble/internal/uuid"
)

type verifyMeta struct {
	nonce   string
	request interface{}
	time    time.Time
}

// Server is a ramble server tasked with storing public keys, encrypted
// messages, and hello-verify handshakes.
type Server struct {
	dur    time.Duration
	active map[string]verifyMeta
	public *splay.Splay

	mu sync.Mutex
}

// NewServer creates a new server. dur is the duration that hello-verify
// handshakes may remain active. publicDir is the directory public keys will be
// stored in.
func NewServer(dur time.Duration, publicDir string) (ret *Server, err error) {
	ret = &Server{
		dur:    dur,
		active: make(map[string]verifyMeta),
	}

	ret.public, err = splay.NewSplay(publicDir, 2)

	if err == nil {
		go ret.prune()
	}

	return
}

// Used as a globally-persisting goroutine to prune handshakes older than s.dur.
// The handshake time value should still be checked since this cannot remove
// stale handshakes immediately.
func (s *Server) prune() {
	ticker := time.NewTicker(s.dur)

	for {
		select {
		case now := <-ticker.C:
			s.mu.Lock()
			for uuid, m := range s.active {
				if now.UTC().Sub(m.time) > s.dur {
					delete(s.active, uuid)
				}
			}
			s.mu.Unlock()
		}
	}
}

// NewHelloResponse generates a hello response and adds it to the active
// handshake map.
func (s *Server) NewHelloResponse(request interface{}) (*ramble.HelloResponse, error) {
	var h ramble.HelloResponse

	b, err := pgp.NonceHex()

	if err != nil {
		return nil, err
	}

	h.Nonce = string(b)

	h.UUID, err = uuid.UUID()

	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.active[h.UUID]; ok {
		log.Printf("%s -> %s already exists in activeHVs!\n",
			h.UUID, m.time.String())
		return nil, errors.New("the very improbable just happened")
	}

	s.active[h.UUID] = verifyMeta{
		nonce:   h.Nonce,
		request: request,
		time:    time.Now().UTC(),
	}

	return &h, nil
}
