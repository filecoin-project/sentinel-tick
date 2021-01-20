package quotetracker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func coinbaseproServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/products/FIL-USD/stats" {
			t.Fatal("unexpected path ", path)
		}

		fmt.Fprintf(w, `
{
  "open": "23.5541",
  "high": "23.6097",
  "low": "22.0638",
  "volume": "214729.266",
  "last": "22.3019",
  "volume_30day": "9045691.422"
}
`)
	}))

	return s
}

func TestCoinbaseproPrice(t *testing.T) {
	s := coinbaseproServer(t)
	defer s.Close()

	cmc := &Coinbasepro{
		url: s.URL,
	}
	q, err := cmc.Price(context.Background(), Pair{Sell: FIL, Buy: USD})
	if err != nil {
		t.Fatal(err)
	}
	if q.Pair.Sell != FIL || q.Pair.Buy != USD {
		t.Error("bad pair set")
	}
	if q.Amount != 22.3019 {
		t.Error("price amount not parsed correctly")
	}
}
