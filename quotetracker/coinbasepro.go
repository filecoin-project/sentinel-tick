package quotetracker

import (
	"context"
	"strconv"
	"time"
)

const coinbaseproURL = "https://api.pro.coinbase.com"

type coinbaseproResponse struct {
	pair Pair

	Last   string `json:"last"`
	Volume string `json:"volume"`
}

func (cbpr *coinbaseproResponse) Quote() (Quote, error) {
	last, err := strconv.ParseFloat(cbpr.Last, 64)
	if err != nil {
		return Quote{}, err
	}
	vol, err := strconv.ParseFloat(cbpr.Volume, 64)
	if err != nil {
		return Quote{}, err
	}

	quote := Quote{
		Pair:      cbpr.pair,
		Timestamp: time.Now(),
		Amount:    last,
		VolumeBase24h: vol,
	}

	return quote, nil
}

type Coinbasepro struct {
	url string // for testing
}

func (ex *Coinbasepro) Price(ctx context.Context, pair Pair) (Quote, error) {
	if ex.url == "" {
		ex.url = coinbaseproURL
	}

	return request(
		ctx,
		ex.url+"/products/"+pair.String()+"/stats",
		nil,
		nil,
		&coinbaseproResponse{pair: pair},
	)
}

func (ex *Coinbasepro) String() string {
	return "coinbasepro"
}
