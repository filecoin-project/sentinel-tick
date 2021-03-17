package quotetracker

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const huobiURL = "https://api.huobi.pro/"

type huobiResponse struct {
	pair Pair

	Status string    `json:"status"`
	Tick   huobiTick `json:"tick"`
	ErrMsg string    `json:"err-msg"`
}

type huobiTick struct {
	Close float64 `json:"close"`
	Vol   float64 `json:"vol"`
}

func (r *huobiResponse) Quote() (Quote, error) {
	if r.Status != "ok" {
		return Quote{}, fmt.Errorf("huobi: bad status: %s", r.ErrMsg)
	}

	quote := Quote{
		Pair:      r.pair,
		Timestamp: time.Now(),
		Amount:    r.Tick.Close,
		VolumeBase24h: r.Tick.Vol,
	}
	return quote, nil
}

var _ Exchange = (*Huobi)(nil)

// Huobi implements fetching quotes from Huobi.
// Note that USD is transparently converted to USDT
type Huobi struct {
	url string
}

// Price fetches the pair information from Huobi.
// For USD, it uses USDT instead, since USD is not traded.
func (ex *Huobi) Price(ctx context.Context, pair Pair) (Quote, error) {
	if ex.url == "" {
		ex.url = huobiURL
	}

	pairOrig := pair
	if pair.Buy == USD {
		pair.Buy = USDT
	}

	return request(
		ctx,
		fmt.Sprintf(
			"%s/market/detail/merged?symbol=%s",
			ex.url, strings.ToLower(
				pair.Sell.Symbol()+pair.Buy.Symbol(),
			),
		),
		nil,
		nil,
		&huobiResponse{pair: pairOrig},
	)
}

func (ex *Huobi) String() string {
	return "huobi"
}
