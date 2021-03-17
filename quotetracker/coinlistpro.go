package quotetracker

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const coinlistproURL = "https://trade-api.coinlist.co"

type coinlistproResponse struct {
	pair Pair

	Message       string               `json:"message"`
	LastTrade     coinlistproLastTrade `json:"last_trade"`
	VolumeBase24h string               `json:"volume_base_24h"`
}

type coinlistproLastTrade struct {
	Price       string    `json:"price"`
	LogicalTime time.Time `json:"logical_time"`
}

func (r *coinlistproResponse) Quote() (Quote, error) {
	if r.Message != "" {
		return Quote{}, fmt.Errorf("coinlistpro: error: %s", r.Message)
	}

	price, err := strconv.ParseFloat(r.LastTrade.Price, 64)
	if err != nil {
		return Quote{}, fmt.Errorf("coinlistpro: error parsing price: %w", err)
	}

	vol, err := strconv.ParseFloat(r.VolumeBase24h, 64)
	if err != nil {
		return Quote{}, fmt.Errorf("coinlistpro: error parsing volume: %w", err)
	}

	quote := Quote{
		Pair: r.pair,
		//Timestamp: r.LastTrade.LogicalTime,
		// not much trading going on. Prefer this to having blanks.
		Timestamp:     time.Now(),
		Amount:        price,
		VolumeBase24h: vol,
	}
	return quote, nil
}

type Coinlistpro struct {
	url string
}

func (ex *Coinlistpro) Price(ctx context.Context, pair Pair) (Quote, error) {
	if ex.url == "" {
		ex.url = coinlistproURL
	}

	return request(
		ctx,
		ex.url+"/v1/symbols/"+pair.String()+"/summary",
		nil,
		nil,
		&coinlistproResponse{pair: pair},
	)
}

func (ex *Coinlistpro) String() string {
	return "coinlistpro"
}
