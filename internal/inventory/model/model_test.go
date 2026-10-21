package model

import (
	"errors"
	"testing"
)

func TestCanonicalOperationIdentity(t *testing.T) {
	base := Request{DatasetID: "dataset", OrderID: "order", OperationID: "key", Method: "reserve", Items: []Item{{"B", 2}, {"A", 1}}}
	normalized, hash, err := Canonical(base)
	if err != nil {
		t.Fatal(err)
	}
	if base.Items[0].SKU != "B" || normalized.Items[0].SKU != "A" {
		t.Fatal("canonicalization mutates caller or does not sort")
	}
	reordered := base
	reordered.Items = []Item{{"A", 1}, {"B", 2}}
	if _, same, e := Canonical(reordered); e != nil || hash != same {
		t.Fatal("SKU order changed operation identity", e)
	}
	changed := base
	changed.OrderID = "other"
	if _, different, e := Canonical(changed); e != nil || hash == different {
		t.Fatal("different order reused request identity", e)
	}
}

func TestInvalidRequests(t *testing.T) {
	for _, r := range []Request{
		{},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "reserve"},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "reserve", Items: []Item{{"A", 0}}},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "reserve", Items: []Item{{"A", 1}, {"A", 2}}},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "commit", Items: []Item{{"A", 1}}},
		{DatasetID: "../secret", OrderID: "o", OperationID: "k", Method: "release"},
	} {
		if _, _, err := Canonical(r); !errors.Is(err, ErrInvalid) {
			t.Fatalf("request accepted: %+v", r)
		}
	}
}
