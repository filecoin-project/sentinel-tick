package quotetracker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func krakenServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pair := r.URL.Query().Get("pair")
		if pair == "" {
			t.Fatal("pair not set")
		}

		fmt.Fprintf(w, `
{
  "error": [],
  "result": {
    "%s": {
      "a": [
        "21.59500",
        "50",
        "50.000"
      ],
      "b": [
        "21.57000",
        "232",
        "232.000"
      ],
      "c": [
        "21.58400",
        "11.88051334"
      ],
      "v": [
        "23559.48334637",
        "48009.68384317"
      ],
      "p": [
        "21.49513",
        "21.40545"
      ],
      "t": [
        5287,
        6085
      ],
      "l": [
        "21.25100",
        "20.92200"
      ],
      "h": [
        "21.89900",
        "21.89900"
      ],
      "o": "21.63600"
    }
  }
}
`, pair)
	}))

	return s
}

func TestKrakenPrice(t *testing.T) {
	s := krakenServer(t)
	defer s.Close()

	cmc := &Kraken{
		url: s.URL + "/0/public/Ticker",
	}
	q, err := cmc.Price(context.Background(), Pair{Sell: FIL, Buy: USD})
	if err != nil {
		t.Fatal(err)
	}
	if q.Pair.Sell != FIL || q.Pair.Buy != USD {
		t.Error("bad pair set")
	}
	if q.Amount != 21.58400 {
		t.Error("price amount not parsed correctly")
	}
}
