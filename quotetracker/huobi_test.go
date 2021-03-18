package quotetracker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func huobiServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/market/detail/merged" {
			t.Fatal("unexpected path ", path)
		}
		if s := r.URL.Query().Get("symbol"); s != "filusdt" {
			t.Fatal("unexpected symbol ", s)
		}

		fmt.Fprintf(w, `
{
  "status": "ok",
  "ch": "market.filusdt.detail.merged",
  "ts": 1611843696655,
  "tick": {
    "amount": 1020192.7761668526,
    "open": 21.4244,
    "close": 21.9954,
    "high": 22.0688,
    "id": 201447395584,
    "count": 83349,
    "low": 21.21,
    "version": 201447395584,
    "ask": [
      21.9953,
      2.2427
    ],
    "vol": 22051370.789555945,
    "bid": [
      21.9939,
      0.5554
    ]
  }
}
`)
	}))

	return s
}

func TestHuobiPrice(t *testing.T) {
	s := huobiServer(t)
	defer s.Close()

	h := &Huobi{
		url: s.URL,
	}
	q, err := h.Price(context.Background(), Pair{Sell: FIL, Buy: USD})
	if err != nil {
		t.Fatal(err)
	}
	if q.Pair.Sell != FIL || q.Pair.Buy != USD {
		t.Error("bad pair set")
	}
	if q.Amount != 21.9954 {
		t.Error("price amount not parsed correctly", q.Amount)
	}

	if q.VolumeBase24h != 1020192.7761668526 {
		t.Error("volume amount not parsed correctly", q.VolumeBase24h)
	}
}
