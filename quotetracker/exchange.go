package quotetracker

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Currency represents a Market currency
type Currency int

// Supported currencies to obtain quotes in.
const (
	EUR Currency = iota
	USD
	FIL
)

// Symbol returns the symbol for a currency.
func (cur Currency) Symbol() string {
	switch cur {
	case EUR:
		return "EUR"
	case USD:
		return "USD"
	case FIL:
		return "FIL"
	default:
		return "UNKNOWN"
	}
}

// CurrencyFromSymbol returns
func CurrencyFromSymbol(symbol string) (Currency, error) {
	var c Currency
	for c = 0; c.Symbol() != Currency(-1).Symbol(); c++ {
		if c.Symbol() == symbol {
			return c, nil
		}
	}
	return -1, errors.New("unsupported currency")
}

// Pair represents two currencies which are exchanged.
type Pair struct {
	Sell Currency
	Buy  Currency
}

// Quote provides price information for a given pair.
type Quote struct {
	Pair      Pair
	Timestamp time.Time
	Amount    float64
}

func (q Quote) String() string {
	return fmt.Sprintf("%s-%s: %f (%s)", q.Pair.Sell.Symbol(), q.Pair.Buy.Symbol(), q.Amount, q.Timestamp)
}

// An exchange returns current Quotes for a given pair.
type Exchange interface {
	Price(context.Context, Pair) (Quote, error)
	String() string
}
