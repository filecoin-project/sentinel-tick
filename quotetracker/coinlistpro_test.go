package quotetracker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func coinlistproServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/v1/symbols/FIL-USD/summary" {
			t.Fatal("unexpected path ", path)
		}

		fmt.Fprintf(w, `
{
  "type": "spot",
  "last_price": "76.35000000",
  "lowest_ask": "76.55000000",
  "highest_bid": "76.14000000",
  "last_trade": {
    "price": "75.41000000",
    "volume": "481.0000",
    "imbalance": "61.6312",
    "logicalTime": "2021-03-17T15:36:13.000Z",
    "auctionCode": "FIL-USD-2021-03-17T15:36:13.000Z"
  },
  "volume_base_24h": "66321.3156",
  "volume_quote_24h": "4271801.9941",
  "price_change_percent_24h": "34.30079156",
  "highest_price_24h": "76.28000000",
  "lowest_price_24h": "56.85000000"
}`)
	}))

	return s
}

func TestCoinlistproPrice(t *testing.T) {
	s := coinlistproServer(t)
	defer s.Close()

	cmc := &Coinlistpro{
		url: s.URL,
	}
	q, err := cmc.Price(context.Background(), Pair{Sell: FIL, Buy: USD})
	if err != nil {
		t.Fatal(err)
	}
	if q.Pair.Sell != FIL || q.Pair.Buy != USD {
		t.Error("bad pair set")
	}
	if q.Amount != 75.41 {
		t.Error("price amount not parsed correctly")
	}
	if q.VolumeBase24h != 66321.3156 {
		t.Error("volume amount not parsed correctly")
	}
}
