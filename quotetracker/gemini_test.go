package quotetracker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func geminiServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/v1/pubticker/FILUSD" {
			t.Fatal("unexpected path ", path)
		}

		fmt.Fprintf(w, `
{
  "bid": "22.0282",
  "ask": "22.0514",
  "volume": {
    "FIL": "8340.55228255",
    "USD": "179387.942121277613",
    "timestamp": 1611844200000
  },
  "last": "22.0384"
}
`)
	}))

	return s
}

func TestGeminiPrice(t *testing.T) {
	s := geminiServer(t)
	defer s.Close()

	g := &Gemini{
		url: s.URL,
	}
	q, err := g.Price(context.Background(), Pair{Sell: FIL, Buy: USD})
	if err != nil {
		t.Fatal(err)
	}
	if q.Pair.Sell != FIL || q.Pair.Buy != USD {
		t.Error("bad pair set")
	}
	if q.Amount != 22.0384 {
		t.Error("price amount not parsed correctly")
	}

	if q.VolumeBase24h != 8340.55228255 {
		t.Error("volume amount not parsed correctly")
	}
}
