package quotetracker

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const geminiURL = "https://api.gemini.com"

type geminiResponse struct {
	pair Pair

	Result  string `json:"result"`
	Message string `json:"message"`

	Last   string                 `json:"last"`
	Volume map[string]interface{} `json:"volume"`
}

func (r *geminiResponse) Quote() (Quote, error) {
	if r.Result != "" {
		return Quote{}, fmt.Errorf("gemini error: %s", r.Message)
	}

	last, err := strconv.ParseFloat(r.Last, 64)
	if err != nil {
		return Quote{}, fmt.Errorf("gemini: error parsing amount: %w", err)
	}

	if r.Volume == nil {
		return Quote{}, errors.New("gemini: no volume information")
	}

	volVal, ok := r.Volume[r.pair.Sell.Symbol()]
	if !ok {
		return Quote{}, fmt.Errorf("gemini: no volume information for %s", r.pair.Sell.Symbol())
	}

	volStr, ok := volVal.(string)
	if !ok {
		return Quote{}, fmt.Errorf("gemini: bad volume information for %s: %+v", r.pair.Sell.Symbol(), volVal)
	}

	vol, err := strconv.ParseFloat(volStr, 64)
	if err != nil {
		return Quote{}, fmt.Errorf("gemini: error parsing volume: %w", err)
	}

	quote := Quote{
		Pair:      r.pair,
		Timestamp: time.Now(),
		Amount:    last,
		VolumeBase24h: vol,
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
