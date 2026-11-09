package model

import (
	"errors"
	"reflect"
	"testing"
)

func TestCanonicalIdentity(t *testing.T) {
	r := Request{DatasetID: "dataset", RequestKey: "key", Items: []Item{{SKU: "b", Quantity: 2}, {SKU: "a", Quantity: 1}}}
	before := append([]Item(nil), r.Items...)
	sorted, hash, err := Canonical(r)
	if err != nil || !reflect.DeepEqual(r.Items, before) {
		t.Fatal("validation mutated input", err)
	}
	_, again, err := Canonical(sorted)
	if err != nil || hash != again {
		t.Fatal("ordering changed identity")
	}
	sorted.Items[0].Quantity++
	_, changed, _ := Canonical(sorted)
	if hash == changed {
		t.Fatal("quantity not in identity")
	}
	r.Items = append(r.Items, Item{SKU: "a", Quantity: 1})
	if _, _, err = Canonical(r); !errors.Is(err, ErrInvalid) {
		t.Fatal("duplicate SKU accepted")
	}
	r.Items = []Item{{SKU: "a", Quantity: 0}}
	if _, _, err = Canonical(r); !errors.Is(err, ErrInvalid) {
		t.Fatal("zero quantity accepted")
	}
}
