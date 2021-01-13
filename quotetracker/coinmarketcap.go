package quotetracker

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"
)

const coinMarketCapURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"

type coinMarketCapResponse struct {
	pair Pair

	Data   map[string]*coinMarketCapData `json:"data"`
	Status coinMarketCapStatus           `json:"status"`
}

type coinMarketCapStatus struct {
	ErrorMessage string    `json:"error_message"`
	Timestamp    time.Time `json:"timestamp"`
}

type coinMarketCapData struct {
	Quote map[string]*coinMarketCapQuote `json:"quote"`
}

type coinMarketCapQuote struct {
	Price     float64 `json:"price"`
	Volume24h float64 `json:"volume_24h"`
}

func (r *coinMarketCapResponse) Quote() (Quote, error) {
	if apiErr := r.Status.ErrorMessage; apiErr != "" {
		return Quote{}, errors.New(apiErr)
	}

	sell := r.pair.Sell.Symbol()
	buy := r.pair.Buy.Symbol()

	if r.Data == nil ||
		r.Data[sell] == nil ||
		r.Data[sell].Quote == nil ||
		r.Data[sell].Quote[buy] == nil {
		return Quote{}, errors.New("expected data not found in response")
	}

	quote := Quote{
		Pair:      r.pair,
		Timestamp: r.Status.Timestamp,
		Amount:    r.Data[sell].Quote[buy].Price,
	}
	return quote, nil
}

var _ Exchange = (*CoinMarketCap)(nil)

// c9bd50b1-7a2b-4b5c-9ede-b12a949cc96b

// CoinMarketCap implements fetching from coinmarketcap.com.
type CoinMarketCap struct {
	Token string
	TTL   time.Duration

	values map[Pair]Quote
	url    string // for testing
}

// Price obtains the current Quote for the given pair. It returns cached values
// if a previous Quote is within its TTL.
func (ex *CoinMarketCap) Price(ctx context.Context, pair Pair) (Quote, error) {
	if ex.values == nil {
		ex.values = make(map[Pair]Quote)
	}

	if ex.url == "" {
		ex.url = coinMarketCapURL
	}

	cached, ok := ex.values[pair]
	if ok && time.Since(cached.Timestamp) <= ex.TTL {
		return cached, nil
	}

	headers := make(http.Header)
	headers.Set("Accepts", "application/json")
	headers.Add("X-CMC_PRO_API_KEY", ex.Token)

	q := url.Values{}
	q.Add("symbol", pair.Sell.Symbol())
	q.Add("convert", pair.Buy.Symbol())

	quote, err := request(
		ctx,
		ex.url,
		q,
		headers,
		&coinMarketCapResponse{pair: pair},
	)
	ex.values[pair] = quote
	return quote, err
}

func (ex *CoinMarketCap) String() string {
	return "coinmarketcap"
}
