package quotetracker

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const geminiURL = "https://api.gemini.com"

type geminiResponse struct {
	pair Pair

	Result  string `json:"result"`
	Message string `json:"message"`

	Last string `json:"last"`
}

func (r *geminiResponse) Quote() (Quote, error) {
	if r.Result != "" {
		return Quote{}, fmt.Errorf("gemini error: %s", r.Message)
	}

	v, err := strconv.ParseFloat(r.Last, 64)
	if err != nil {
		return Quote{}, fmt.Errorf("gemini: error parsing amount: %w", err)
	}

	quote := Quote{
		Pair:      r.pair,
		Timestamp: time.Now(),
		Amount:    v,
	}
	return quote, nil
}

var _ Exchange = (*Gemini)(nil)

// Gemini fetches price information from the Gemini exchange.
type Gemini struct {
	url string
}

// Price fetches the last ticker price from Gemini.
func (ex *Gemini) Price(ctx context.Context, pair Pair) (Quote, error) {
	if ex.url == "" {
		ex.url = geminiURL
	}

	return request(
		ctx,
		ex.url+"/v1/pubticker/"+pair.Sell.Symbol()+pair.Buy.Symbol(),
		nil,
		nil,
		&geminiResponse{pair: pair},
	)
}

func (ex *Gemini) String() string {
	return "gemini"
}
