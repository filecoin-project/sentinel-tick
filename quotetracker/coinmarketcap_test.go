package quotetracker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func coinMarketCapServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("X-CMC_PRO_API_KEY"); token == "" {
			t.Fatal("token not set")
		}

		sell := r.URL.Query().Get("symbol")
		convert := r.URL.Query().Get("convert")
		if sell == "" || convert == "" {
			t.Fatal("pair not set")
		}
		buys := strings.Split(convert, ",")
		if len(buys) != 2 {
			t.Fatal("wrong convert argument")
		}

		fmt.Fprintf(w, `
{
  "status": {
    "timestamp": "2021-01-12T13:13:48.439Z",
    "error_code": 0,
    "error_message": null,
    "elapsed": 16,
    "credit_count": 1,
    "notice": null
  },
  "data": {
    "%s": {
      "id": 2280,
      "name": "Filecoin",
      "symbol": "%s",
      "slug": "filecoin",
      "num_market_pairs": 80,
      "date_added": "2017-12-13T00:00:00.000Z",
      "tags": [
        "mineable",
        "distributed-computing",
        "filesharing",
        "ipfs"
      ],
      "max_supply": 2000000000,
      "circulating_supply": 44584205,
      "total_supply": 44584205,
      "is_active": 1,
      "platform": null,
      "cmc_rank": 39,
      "is_fiat": 0,
      "last_updated": "2021-01-12T13:13:03.000Z",
      "quote": {
        "%s": {
          "price": 1,
          "volume_24h": 63993291.66774165,
          "percent_change_1h": 0,
          "percent_change_24h": 0,
          "percent_change_7d": 0,
          "percent_change_30d": 0,
          "percent_change_60d": 0,
          "percent_change_90d": 0,
          "market_cap": 58362381.00000001,
          "last_updated": "2021-03-18T15:06:04.000Z"
        },
        "%s": {
          "price": 21.58619654190641,
          "volume_24h": 251752867.90136266,
          "percent_change_1h": -0.45896058,
          "percent_change_24h": -0.91329396,
          "percent_change_7d": -0.88271718,
          "market_cap": 962403411.7946465,
          "last_updated": "2021-01-12T13:13:03.000Z"
        }
      }
    }
  }
}
`, sell, sell, buys[0], buys[1])
	}))

	return s
}

func TestCoinMarketCapPrice(t *testing.T) {
	s := coinMarketCapServer(t)
	defer s.Close()

	cmc := &CoinMarketCap{
		Token: "auth",
		TTL:   time.Second,
		url:   s.URL,
	}
	q, err := cmc.Price(context.Background(), Pair{Sell: FIL, Buy: USD})
	if err != nil {
		t.Fatal(err)
	}
	if q.Pair.Sell != FIL || q.Pair.Buy != USD {
		t.Error("bad pair set")
	}
	if q.Amount != 21.58619654190641 {
		t.Error("price amount not parsed correctly")
	}

	if q.VolumeBase24h != 63993291.66774165 {
		t.Error("volume amount not parsed correctly")
	}
}
