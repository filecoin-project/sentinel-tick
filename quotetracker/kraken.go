package quotetracker

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const krakenURL = "https://api.kraken.com"

type krakenResponse struct {
	pair Pair

	Error  []interface{}            `json:"error"`
	Result map[string]*krakenResult `json:"result"`
}

type krakenResult struct {
	LastTradeClosed []string `json:"c"`
	Volume          []string `json:"v"`
}

func (r *krakenResponse) Quote() (Quote, error) {
	if len(r.Error) > 0 {
		return Quote{}, fmt.Errorf("kraken: response has errors: %s", r.Error)
	}

	// kraken returns pair strings like XXBTZUSD instead of BTCUSD
	// Since we request a single one, let's just assume what we requested
	// comes back.
	var pairStr string
	for k := range r.Result {
		pairStr = k
		break
	}

	if r.Result == nil ||
		r.Result[pairStr] == nil ||
		len(r.Result[pairStr].LastTradeClosed) < 2 ||
		len(r.Result[pairStr].Volume) < 2 {
		return Quote{}, fmt.Errorf("kraken: unexpected response: %+v", r)
	}

	price, err := strconv.ParseFloat(r.Result[pairStr].LastTradeClosed[0], 64)
	if err != nil {
		return Quote{}, fmt.Errorf("kraken: error parsing price: %w", err)
	}

	vol, err := strconv.ParseFloat(r.Result[pairStr].Volume[1], 64)
	if err != nil {
		return Quote{}, fmt.Errorf("kraken: error parsing volume: %w", err)
	}

	quote := Quote{
		Pair:          r.pair,
		Timestamp:     time.Now(),
		Amount:        price,
		VolumeBase24h: vol,
	}

	return quote, nil
}

var _ Exchange = (*Kraken)(nil)

// Kraken implements fetching quotes from Kraken.
type Kraken struct {
	url string // for testing
}

// Price fetches the last trade price from Kraken.
func (ex *Kraken) Price(ctx context.Context, pair Pair) (Quote, error) {
	if ex.url == "" {
		ex.url = krakenURL
	}

	sell := pair.Sell.Symbol()
	buy := pair.Buy.Symbol()
	if sell == "BTC" {
		sell = "XBT"
	}
	if buy == "BTC" {
		buy = "XBT"
	}

	q := url.Values{}
	q.Add("pair", sell+buy)
	return request(
		ctx,
		ex.url+"/0/public/Ticker",
		q,
		nil,
		&krakenResponse{pair: pair},
	)
}

func (ex *Kraken) String() string {
	return "kraken"
}
