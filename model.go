package main

import (
	"context"

	"github.com/filecoin-project/sentinel-tick/quotetracker"
	pg "github.com/go-pg/pg/v10"
)

const schema = "currencies"

// Quote stores FIL price information.
type Quote struct {
	//lint:ignore U1000 hit for go-pg
	tableName     struct{} `pg:"currencies.quotes"`
	Height        int64    `pg:",pk,notnull"`
	Price         int64    `pg:",notnull"`
	VolumeBase24h int64    `pg:"volume_base_24h,notnull"`
	Exchange      string   `pg:",pk,notnull"`
	Sell          string   `pg:",pk,notnull"`
	Buy           string   `pg:",pk,notnull"`
}

// NewQuote creates a new FIL quote for the database.
func NewQuote(h int64, ex string, q quotetracker.Quote) Quote {
	return Quote{
		Height:        h,
		Price:         toMicro(q.Amount),
		VolumeBase24h: toMicro(q.VolumeBase24h),
		Exchange:      ex,
		Sell:          q.Pair.Sell.Symbol(),
		Buy:           q.Pair.Buy.Symbol(),
	}
}

// toMicro converts a float unit to int micro-units for storage.
func toMicro(amount float64) int64 {
	return int64(amount * 1000000) // micro
}

// func fromMicro(amount int64) float64 {
// 	return float64(amount / 1000000) // micro
// }

// Persist uses a transaction to insert a quote in the DB.
func (q Quote) Persist(ctx context.Context, tx *pg.Tx) error {
	_, err := tx.ModelContext(ctx, &q).
		OnConflict("do nothing").
		Insert()
	return err
}

// Quotes groups multiple Quotes so they can be persisted together.
type Quotes []Quote

// Persist uses a transaction to insert multiple quotes in the DB.
func (q Quotes) Persist(ctx context.Context, tx *pg.Tx) error {
	if len(q) == 0 {
		return nil
	}

	_, err := tx.ModelContext(ctx, &q).
		OnConflict("do nothing").
		Insert()
	return err
}
