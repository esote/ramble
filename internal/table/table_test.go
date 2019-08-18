package table

import (
	"testing"

	"github.com/majiru/ramble/internal/uuid"
)

func TestTable(t *testing.T) {
	table, err := NewTable("test_table", 2, uuid.LenUUID)

	if err != nil {
		t.Fatal(err)
	}

	const (
		n   = 10
		key = "9baacc8baed73d1f115d10d069a4ee63i"
	)

	rows := make([]string, n)

	for i := 0; i < n; i++ {
		row, err := uuid.UUID()

		if err != nil {
			t.Fatal(err)
		}

		if err = table.Insert(key, row); err != nil {
			t.Fatal(err)
		}

		rows[i] = row
	}

	// Check Index matches the expected list of messages. Index returns msgs
	// in the order of newest to oldest, so carefully iterate and slice.
	for take := 0; take <= n; take++ {
		index, err := table.Index(key, uint64(take))

		if err != nil {
			t.Fatal(err)
		}

		for i, m := range rows[len(rows)-take:] {
			if m != index[len(index)-i-1] {
				t.Fatalf("mismatch at take: %d", take)
			}
		}
	}

	if err = table.Splay.RemoveAll(); err != nil {
		t.Fatal(err)
	}
}

func TestTableUnique(t *testing.T) {
	table, err := NewTable("test_table", 2, 1)

	if err != nil {
		t.Fatal(err)
	}

	const (
		key = "9baacc8baed73d1f115d10d069a4ee63i"
	)

	rows := []string{
		"a",
		"b",
		"c",
		"d",
	}

	insert := []string{
		"a",
		"b",
		"c",
		"b",
		"d",
		"a",
		"d",
	}

	for _, row := range insert {
		if err = table.InsertUnique(key, row); err != nil {
			t.Fatal(err)
		}
	}

	index, err := table.Index(key, 0)

	if err != nil {
		t.Fatal(err)
	}

	if len(index) != len(rows) {
		t.Fatal("incorrect length")
	}

	for i, row := range index {
		if rows[len(rows)-i-1] != row {
			t.Fatalf("row %d mismatch", i)
		}
	}

	if err = table.Splay.RemoveAll(); err != nil {
		t.Fatal(err)
	}
}
