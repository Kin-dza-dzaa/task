// Package entity implements entities layer of the app.
// Describes entities and some errors that are used in all layers.
package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Currency rate to KZT.
// JSON tags used only for serialization.
type Rate struct {
	Quant   uint            `json:"quant"`
	Title   string          `json:"title"`
	Code    string          `json:"code"`
	Change  string          `json:"change"`
	Rate    decimal.Decimal `json:"rate"`
	ValidAt time.Time       `json:"valid_at"`
}
