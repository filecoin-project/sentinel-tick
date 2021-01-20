package quotetracker

import (
	"context"
	"strconv"
	"time"
)

const coinlistproURL = "https://trade-api.coinlist.co"

type coinlistproResponse struct {
	pair Pair

	LastTrade coinlistproLastTrade `json:"last_trade"`
}

type coinlistproLastTrade struct {
	Price       string    `json:"price"`
	LogicalTime time.Time `json:"logical_time"`
}

func (r *coinlistproResponse) Quote() (Quote, error) {
	v, err := strconv.ParseFloat(r.LastTrade.Price, 64)
	if err != nil {
		return Quote{}, err
	}

	quote := Quote{
		Pair: r.pair,
		//Timestamp: r.LastTrade.LogicalTime,
		// not much trading going on. Prefer this to having blanks.
		Timestamp: time.Now(),
		Amount:    v,
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
		ex.url+"/v1/symbols/"+pair.String()+"/quote",
		nil,
		nil,
		&coinlistproResponse{pair: pair},
	)
}

func (ex *Coinlistpro) String() string {
	return "coinlistpro"
}
