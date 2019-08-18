package cindex

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/esote/util/splay"
	"github.com/majiru/ramble/internal/uuid"
)

// CIndex handles ramble's conversation indexing and insertion. CIndex is a
// wrapper around Splay. Insertion is constant-time, indexing scales linearly
// with the want argument.
type CIndex struct {
	Splay *splay.Splay
}

// NewCIndex creates a new CIndex.
func NewCIndex(convo string, cutoff uint64) (*CIndex, error) {
	ret := &CIndex{}

	var err error

	ret.Splay, err = splay.NewSplay(convo, cutoff)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Index list of latest message UUIDs. If want == 0, all messages will be
// returned.
func (c *CIndex) Index(convo string, want uint64) ([]string, error) {
	if !c.Splay.Exists(convo) {
		return nil, errors.New("no such conversation")
	}

	f, err := c.Splay.Open(convo)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	count, err := readCount(f)

	if err != nil {
		return nil, err
	}

	var items []string

	if want == 0 || want >= count {
		items = make([]string, count)
	} else {
		items = make([]string, want)
	}

	offset := int64(len(items) * uuid.LenUUID)

	if _, err = f.Seek(-offset, io.SeekEnd); err != nil {
		return nil, err
	}

	buf := make([]byte, uuid.LenUUID)

	for i := len(items) - 1; i >= 0; i-- {
		if _, err = f.Read(buf); err != nil {
			return nil, err
		}

		items[i] = string(buf)
	}

	return items, nil
}

// Insert a new message into conversation index.
func (c *CIndex) Insert(convo string) (string, error) {
	msg, err := uuid.UUID()

	if err != nil {
		return "", err
	}

	if !c.Splay.Exists(convo) {
		return msg, c.createFile(convo, msg)
	}

	return msg, c.appendFile(convo, msg)
}

func (c *CIndex) createFile(convo, msg string) error {
	var b bytes.Buffer
	b.Grow(64 + uuid.LenUUID)

	_, _ = b.Write(encodeCount(1))
	_, _ = b.WriteString(msg)

	return c.Splay.Write(convo, b.Bytes())
}

func (c *CIndex) appendFile(convo, msg string) error {
	f, err := c.Splay.OpenFile(convo, os.O_RDWR, 0)

	if err != nil {
		return err
	}

	defer f.Close()

	count, err := readCount(f)

	if err != nil {
		return err
	}

	if _, err = f.WriteAt(encodeCount(count+1), 0); err != nil {
		return err
	}

	if _, err = f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	if _, err = f.WriteString(msg); err != nil {
		return err
	}

	return nil
}

func encodeCount(count uint64) []byte {
	return []byte(fmt.Sprintf("%064b", count))
}

// Read and decode count. Assumes seek offset 0.
func readCount(f *os.File) (out uint64, err error) {
	buf := make([]byte, 64)

	if _, err = f.Read(buf); err != nil {
		return
	}

	_, err = fmt.Sscanf(string(buf), "%064b", &out)
	return
}
