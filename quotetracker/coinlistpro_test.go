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
		if path != "/v1/symbols/FIL-USD/quote" {
			t.Fatal("unexpected path ", path)
		}

		fmt.Fprintf(w, `
{
  "last_trade": {
    "price": "22.30000000",
    "volume": "9.7465",
    "imbalance": "190.2535",
    "logical_time": "2021-01-20T12:39:33.000Z",
    "auction_code": "FIL-USD-2021-01-20T12:39:33.000Z"
  },
  "quote": {
    "ask": "22.30000000",
    "ask_size": "312.1919",
    "bid": "22.25000000",
    "bid_size": "312.0856"
  },
  "after_auction_code": "FIL-USD-2021-01-20T12:43:09.000Z",
  "call_time": "2021-01-20T12:43:09.086Z",
  "logical_time": "2021-01-20T12:43:09.000Z"
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
	if q.Amount != 22.30 {
		t.Error("price amount not parsed correctly")
	}
}
