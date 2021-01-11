package quotetracker

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const coinMarketCapURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"

type apiData struct {
	Data   map[string]*curData `json:"data"`
	Status status              `json:"status"`
}

type status struct {
	ErrorMessage string    `json:"error_message"`
	Timestamp    time.Time `json:"timestamp"`
}

type curData struct {
	Quote map[string]*quoteData `json:"quote"`
}

type quoteData struct {
	Price     float64 `json:price`
	Volume24h float64 `json:volume_24h`
}

// c9bd50b1-7a2b-4b5c-9ede-b12a949cc96b

// CoinMarketCap implements fetching from coinmarketcap.com.
type CoinMarketCap struct {
	Token string
	TTL   time.Duration

	values map[Pair]Quote
	client *http.Client
}

// Price obtains the current Quote for the given pair. It returns cached values
// if a previous Quote is within its TTL.
func (ex *CoinMarketCap) Price(ctx context.Context, pair Pair) (Quote, error) {
	if ex.client == nil {
		ex.client = &http.Client{}
	}

	if ex.values == nil {
		ex.values = make(map[Pair]Quote)
	}

	cached, ok := ex.values[pair]
	if ok && time.Since(cached.Timestamp) <= ex.TTL {
		return cached, nil
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		coinMarketCapURL,
		nil,
	)

	if err != nil {
		return Quote{}, err
	}

	q := url.Values{}
	q.Add("symbol", pair.Sell.Symbol())
	q.Add("convert", pair.Buy.Symbol())

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", ex.Token)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return Quote{}, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Quote{}, err
	}
	var parsedResp apiData
	err = json.Unmarshal(respBody, &parsedResp)
	if err != nil {
		return Quote{}, err
	}

	if apiErr := parsedResp.Status.ErrorMessage; apiErr != "" {
		return Quote{}, errors.New(apiErr)
	}

	if parsedResp.Data == nil ||
		parsedResp.Data[pair.Sell.Symbol()] == nil ||
		parsedResp.Data[pair.Sell.Symbol()].Quote == nil ||
		parsedResp.Data[pair.Sell.Symbol()].Quote[pair.Buy.Symbol()] == nil {
		return Quote{}, errors.New("expected data not found in response")
	}

	quote := Quote{
		Pair:      pair,
		Timestamp: parsedResp.Status.Timestamp,
		Amount:    parsedResp.Data[pair.Sell.Symbol()].Quote[pair.Buy.Symbol()].Price,
	}

	ex.values[pair] = quote
	return quote, nil
}

func (ex *CoinMarketCap) String() string {
	return "coinmarketcap"
}
