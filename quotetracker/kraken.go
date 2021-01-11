package quotetracker

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const krakenURL = "https://api.kraken.com/0/public/Ticker"

type krakenResponse struct {
	pair Pair

	Error  []interface{}            `json:"error"`
	Result map[string]*krakenResult `json:"result"`
}

type krakenResult struct {
	LastTradeClosed []string `json:"c"`
}

func (r *krakenResponse) Quote() (Quote, error) {
	if len(r.Error) > 0 {
		return Quote{}, fmt.Errorf("response has errors: %s", r.Error)
	}

	pairStr := r.pair.Sell.Symbol() + r.pair.Buy.Symbol()

	if r.Result == nil ||
		r.Result[pairStr] == nil ||
		len(r.Result[pairStr].LastTradeClosed) < 2 {
		return Quote{}, fmt.Errorf("unexpected response: %+v", r)
	}

	price, err := strconv.ParseFloat(r.Result[pairStr].LastTradeClosed[0], 64)
	if err != nil {
		return Quote{}, err
	}

	quote := Quote{
		Pair:      r.pair,
		Timestamp: time.Now(),
		Amount:    price,
	}

	return quote, nil
}

// Kraken implements fetching quotes from Kraken.
type Kraken struct {
}

// Price fetches the last trade price from Kraken.
func (ex *Kraken) Price(ctx context.Context, pair Pair) (Quote, error) {
	q := url.Values{}
	q.Add("pair", pair.Sell.Symbol()+pair.Buy.Symbol())
	return request(
		ctx,
		krakenURL,
		q,
		nil,
		&krakenResponse{pair: pair},
	)
}

func (ex *Kraken) String() string {
	return "kraken"
}
