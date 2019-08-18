package cindex

import (
	"testing"
)

func TestCIndex(t *testing.T) {
	c, err := NewCIndex("test_convo", 2)

	if err != nil {
		t.Fatal(err)
	}

	const (
		n    = 10
		uuid = "9baacc8baed73d1f115d10d069a4ee63i"
	)

	msgs := make([]string, n)

	for i := 0; i < n; i++ {
		msg, err := c.Insert(uuid)

		if err != nil {
			t.Fatal(err)
		}

		msgs[i] = msg
	}

	// Check Index matches the expected list of messages. Index returns msgs
	// in the order of newest to oldest, so carefully iterate and slice.
	for take := 0; take <= n; take++ {
		index, err := c.Index(uuid, uint64(take))

		if err != nil {
			t.Fatal(err)
		}

		for i, m := range msgs[len(msgs)-take : len(msgs)] {
			if m != index[len(index)-i-1] {
				t.Fatalf("mismatch at take: %d", take)
			}
		}
	}

	if err = c.Splay.RemoveAll(); err != nil {
		t.Fatal(err)
	}
}
