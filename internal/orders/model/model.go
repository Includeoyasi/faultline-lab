// Package model defines Orders' request identity and durable reserve step.
package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"slices"
	"strings"
)

var (
	ErrInvalid     = errors.New("invalid orders request")
	ErrKeyConflict = errors.New("request key already has another payload")
	ErrDataset     = errors.New("dataset does not match orders")
	ErrNotFound    = errors.New("order not found")
	ErrStale       = errors.New("reserve attempt no longer owns the job")
	identifier     = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)
	sku            = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`)
)

const (
	Accepted       = "ACCEPTED"
	Reserving      = "RESERVING"
	PaymentPending = "PAYMENT_PENDING"
	Rejected       = "REJECTED"
)

type Item struct {
	SKU      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}
type PricedItem struct {
	SKU            string `json:"sku"`
	Quantity       int64  `json:"quantity"`
	UnitPriceMinor int64  `json:"unit_price_minor"`
}
type Request struct {
	DatasetID  string `json:"dataset_id"`
	RequestKey string `json:"request_key"`
	Items      []Item `json:"items"`
}
type Acceptance struct {
	DatasetID  string
	OrderID    string
	RequestKey string
}
type Reservation struct {
	Version    int64        `json:"version"`
	Items      []PricedItem `json:"items"`
	TotalMinor int64        `json:"total_minor"`
	Currency   string       `json:"currency"`
}
type Order struct {
	DatasetID     string       `json:"dataset_id"`
	OrderID       string       `json:"order_id"`
	RequestKey    string       `json:"request_key"`
	OperationID   string       `json:"reserve_operation_id"`
	State         string       `json:"state"`
	Version       int64        `json:"version"`
	Items         []Item       `json:"items"`
	Reservation   *Reservation `json:"reservation,omitempty"`
	Rejection     string       `json:"rejection,omitempty"`
	BlockedReason string       `json:"blocked_reason,omitempty"`
}
type Job struct {
	Order        Order
	Generation   int64
	OrderVersion int64
	Attempts     int64
}

// Resolution is one verified peer outcome or a bounded diagnostic code.
type Resolution struct {
	Reservation   *Reservation
	Rejection     string
	BlockedReason string
	RetryReason   string
}

func ValidIdentifier(value string) bool { return identifier.MatchString(value) }
func Canonical(r Request) (Request, string, error) {
	if !ValidIdentifier(r.DatasetID) || !ValidIdentifier(r.RequestKey) || len(r.Items) < 1 || len(r.Items) > 100 {
		return Request{}, "", ErrInvalid
	}
	r.Items = slices.Clone(r.Items)
	slices.SortFunc(r.Items, func(a, b Item) int { return strings.Compare(a.SKU, b.SKU) })
	for i, item := range r.Items {
		if !sku.MatchString(item.SKU) || item.Quantity <= 0 || i > 0 && item.SKU == r.Items[i-1].SKU {
			return Request{}, "", ErrInvalid
		}
	}
	data, err := json.Marshal(r)
	if err != nil {
		return Request{}, "", err
	}
	hash := sha256.Sum256(data)
	return r, hex.EncodeToString(hash[:]), nil
}
