package table

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/esote/util/splay"
)

// Table handles ramble's tabular indexing and insertion of fixed-length rows.
// Table is a wrapper around Splay.
//
// Time complexities: Index O(want), Insert O(1), InsertUnique O(existing rows).
//
// Space complexities: Insert and InsertUnique O(1), and Index O(want).
type Table struct {
	Splay *splay.Splay

	rowLen int
}

// NewTable creates a new Table.
func NewTable(dir string, cutoff uint64, rowLen int) (t *Table, err error) {
	t = &Table{
		rowLen: rowLen,
	}

	t.Splay, err = splay.NewSplay(dir, cutoff)
	return
}

// Index list of latest table rows. If want == 0, all rows will be returned.
func (t *Table) Index(key string, want uint64) ([]string, error) {
	if !t.Splay.Exists(key) {
		return nil, errors.New("no such conversation")
	}

	f, err := t.Splay.Open(key)

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

	offset := int64(len(items) * t.rowLen)

	if _, err = f.Seek(-offset, io.SeekEnd); err != nil {
		return nil, err
	}

	buf := make([]byte, t.rowLen)

	for i := len(items) - 1; i >= 0; i-- {
		if _, err = f.Read(buf); err != nil {
			return nil, err
		}

		items[i] = string(buf)
	}

	return items, nil
}

// Insert row into table.
func (t *Table) Insert(key, row string) error {
	if len(row) != t.rowLen {
		return errors.New("row length invalid")
	}

	if t.Splay.Exists(key) {
		return t.appendRow(key, row)
	}

	return t.create(key, row)
}

// InsertUnique row into table. If the row already exists, no change is made to
// the key file.
func (t *Table) InsertUnique(key, row string) error {
	if len(row) != t.rowLen {
		return errors.New("row length invalid")
	}

	if t.Splay.Exists(key) {
		return t.appendUniqueRow(key, row)
	}

	return t.create(key, row)
}

func (t *Table) create(key, row string) error {
	var b bytes.Buffer
	b.Grow(64 + t.rowLen)

	_, _ = b.Write(encodeCount(1))
	_, _ = b.WriteString(row)

	return t.Splay.Write(key, b.Bytes())
}

func (t *Table) appendRow(key, row string) error {
	f, err := t.Splay.OpenFile(key, os.O_RDWR, 0)

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

	_, err = f.WriteString(row)
	return err
}

func (t *Table) appendUniqueRow(key, row string) error {
	f, err := t.Splay.OpenFile(key, os.O_RDWR, 0)

	if err != nil {
		return err
	}

	defer f.Close()

	count, err := readCount(f)

	if err != nil {
		return err
	}

	buf := make([]byte, t.rowLen)

	for {
		if _, err = f.Read(buf); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if string(buf) == row {
			return nil
		}
	}

	// At end of file: write row and increment count.

	if _, err = f.WriteString(row); err != nil {
		return err
	}

	_, err = f.WriteAt(encodeCount(count+1), 0)
	return err
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
