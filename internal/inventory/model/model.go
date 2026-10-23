// Package model defines Inventory requests, results and validation shared by storage and transport.
package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
)

var (
	ErrInvalid     = errors.New("invalid inventory request")
	ErrKeyConflict = errors.New("operation key already has another request")
	ErrDataset     = errors.New("dataset does not match inventory")
	ErrNotFound    = errors.New("reservation or product not found")
)

var identifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)
var skuPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`)

const (
	Reserved          = "RESERVED"
	Committed         = "COMMITTED"
	Released          = "RELEASED"
	Applied           = "APPLIED"
	InsufficientStock = "INSUFFICIENT_STOCK"
	UnknownSKU        = "UNKNOWN_SKU"
	InvalidState      = "INVALID_STATE"
	ReservationExists = "RESERVATION_EXISTS"
)

// Item is canonicalized by SKU; duplicate SKUs are rejected.
type Item struct {
	SKU      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}
type PricedItem struct {
	SKU            string `json:"sku"`
	Quantity       int64  `json:"quantity"`
	UnitPriceMinor int64  `json:"unit_price_minor"`
}

// Reservation is the authoritative snapshot when returned by GetReservation.
type Reservation struct {
	DatasetID  string       `json:"dataset_id"`
	OrderID    string       `json:"order_id"`
	State      string       `json:"state"`
	Version    int64        `json:"version"`
	Items      []PricedItem `json:"items"`
	TotalMinor int64        `json:"total_minor"`
	Currency   string       `json:"currency"`
}

// Receipt is immutable, even if the reservation changes after this operation.
type Receipt struct {
	OperationID string       `json:"operation_id"`
	Outcome     string       `json:"outcome"`
	Reservation *Reservation `json:"reservation_at_operation,omitempty"`
}
type Request struct {
	DatasetID   string `json:"dataset_id"`
	OrderID     string `json:"order_id"`
	OperationID string `json:"operation_id"`
	Method      string `json:"method"`
	Items       []Item `json:"items,omitempty"`
}
type Availability struct {
	SKU            string
	Available      int64
	UnitPriceMinor int64
	Currency       string
}

// Canonical validates a request and computes its stable operation identity.
func Canonical(request Request) (Request, string, error) {
	if !identifier.MatchString(request.DatasetID) || !identifier.MatchString(request.OrderID) || !identifier.MatchString(request.OperationID) {
		return Request{}, "", ErrInvalid
	}
	switch request.Method {
	case "reserve":
		if len(request.Items) == 0 || len(request.Items) > 100 {
			return Request{}, "", ErrInvalid
		}
	case "commit", "release":
		if len(request.Items) != 0 {
			return Request{}, "", ErrInvalid
		}
	default:
		return Request{}, "", ErrInvalid
	}
	request.Items = append([]Item(nil), request.Items...)
	sort.Slice(request.Items, func(i, j int) bool { return request.Items[i].SKU < request.Items[j].SKU })
	for i, item := range request.Items {
		if !skuPattern.MatchString(item.SKU) || item.Quantity <= 0 || (i > 0 && item.SKU == request.Items[i-1].SKU) {
			return Request{}, "", ErrInvalid
		}
	}
	data, err := json.Marshal(request)
	if err != nil {
		return Request{}, "", err
	}
	hash := sha256.Sum256(data)
	return request, hex.EncodeToString(hash[:]), nil
}

// ValidIdentifier validates a dataset, order or operation identifier.
func ValidIdentifier(value string) bool { return identifier.MatchString(value) }

// ValidSKU validates a product identifier.
func ValidSKU(value string) bool { return skuPattern.MatchString(value) }
