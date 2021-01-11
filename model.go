package main

import (
	"context"

	"github.com/filecoin-project/sentinel-tick/quotetracker"
	pg "github.com/go-pg/pg/v10"
)

const schema = "filquotes"

// Quote stores FIL price information.
type Quote struct {
	//lint:ignore U1000 hit for go-pg
	tableName struct{} `pg:"filquotes.fil_quotes"`
	Height    int64    `pg:",pk,use_zero,notnull"`
	Price     int64    `pg:",use_zero,notnull"`
	Exchange  string   `pg:",pk,use_zero,notnull`
	Currency  string   `pg:",pk,use_zero,notnull`
}

// NewQuote creates a new FIL quote for the database.
func NewQuote(h int64, ex string, q quotetracker.Quote) Quote {
	return Quote{
		Height:   h,
		Price:    toMicroFIL(q.Amount),
		Exchange: ex,
		Currency: q.Pair.Buy.Symbol(),
	}
}

// toMicroFIL converts FIL price to int for storage.
func toMicroFIL(price float64) int64 {
	return int64(price * 1000000) // microFIL
}

// func fromMicroFIL(price int64) float64 {
// 	return float64(price / 1000000) // microFIL
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
