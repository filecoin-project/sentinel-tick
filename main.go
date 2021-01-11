package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/filecoin-project/sentinel-tick/quotetracker"
	"github.com/go-pg/pg/v10"
	"github.com/urfave/cli/v2"
)

// to be set at build time.
var tag string
var version string = "unset"

func init() {
	if v, err := semver.ParseTolerant(tag); err == nil {
		version = v.String()
	}
}

func main() {

	app := &cli.App{
		Name:    "tick",
		Usage:   "Filecoin Price Monitoring Utility",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "db",
				EnvVars: []string{"SENTINEL_TICK_DB"},
				Value:   "postgres://postgres:password@localhost:5432/postgres?sslmode=disable",
			},
			&cli.IntFlag{
				Name:    "db-pool-size",
				EnvVars: []string{"SENTINEL_TICK_DB_POOL_SIZE"},
				Value:   75,
			},
			&cli.StringFlag{
				Name:    "pairs",
				Usage:   "Comma-separated list of pairs",
				EnvVars: []string{"SENTINEL_TICK_PAIRS"},
				Value:   "FIL-USD",
			},
			&cli.IntFlag{
				Name:    "coinmarketcap",
				Aliases: []string{"cmc"},
				Usage:   "Interval between requests to coinmarketcap. 0 to disable",
				Value:   0,
				EnvVars: []string{"SENTINEL_TICK_COINMARKETCAP"},
			},
			&cli.StringFlag{
				Name:    "coinmarketcap-token",
				Aliases: []string{"cmk-token"},
				Usage:   "API token for CoinMarketCap.com",
				EnvVars: []string{"SENTINEL_TICK_COINMARKETCAP_TOKEN"},
			},
		},
		Action: func(cctx *cli.Context) error {
			db, err := setupDB(
				cctx.Context,
				cctx.String("db"),
				cctx.Int("db-pool-size"),
				"sentinel-tick",
			)
			if err != nil {
				return err
			}
			defer teardownDB(cctx.Context, db)

			if err := createSchema(db); err != nil {
				return err
			}

			exchanges, err := setupExchanges(cctx)
			if err != nil {
				return err
			}

			pairs, err := setupPairs(cctx)
			if err != nil {
				return err
			}

			for {
				select {
				case <-cctx.Context.Done():
					return nil
				default:
				}
				err := fetchQuotes(cctx, db, exchanges, pairs)
				if err != nil {
					return err
				}
				time.Sleep(30 * time.Second)
			}
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Parse flags related to exchanges and initialize them.
func setupExchanges(cctx *cli.Context) ([]quotetracker.Exchange, error) {
	var exchanges []quotetracker.Exchange

	if cmkInterval := cctx.Int("cmc"); cmkInterval > 0 {
		cmc := &quotetracker.CoinMarketCap{
			Token: cctx.String("cmk-token"),
			TTL:   time.Second * time.Duration(cmkInterval),
		}
		exchanges = append(exchanges, cmc)
		fmt.Println("Initialized", cmc)
	}

	return exchanges, nil
}

// Parse the --pairs flag.
func setupPairs(cctx *cli.Context) ([]quotetracker.Pair, error) {
	pairsArr := strings.Split(cctx.String("pairs"), ",")
	pairs := make([]quotetracker.Pair, 0, len(pairsArr))

	for _, pairStr := range pairsArr {
		pairSymbols := strings.Split(pairStr, "-")
		if len(pairSymbols) != 2 {
			return nil, fmt.Errorf("wrong pair: %s", pairStr)
		}

		sell, err := quotetracker.CurrencyFromSymbol(pairSymbols[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", pairSymbols[0], err)
		}
		buy, err := quotetracker.CurrencyFromSymbol(pairSymbols[1])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", pairSymbols[1], err)
		}

		pair := quotetracker.Pair{Sell: sell, Buy: buy}
		pairs = append(pairs, pair)
		fmt.Println("Tracking", pair.String())
	}
	return pairs, nil
}

// For each exchange, for each given pair, fetch quotes and add them to the DB (in a single transaction).
func fetchQuotes(cctx *cli.Context, db *pg.DB, exchanges []quotetracker.Exchange, pairs []quotetracker.Pair) error {
	quotes := make(Quotes, 0, len(exchanges))
	epoch := filEpoch(time.Now())

	fmt.Println("-- Epoch ", epoch)
	for _, ex := range exchanges {
		for _, pair := range pairs {
			q, err := ex.Price(cctx.Context, pair)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(ex, q)
			quotes = append(quotes, NewQuote(
				epoch,
				ex.String(),
				q,
			))
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Close()

	return db.RunInTransaction(cctx.Context, func(tx *pg.Tx) error {
		return quotes.Persist(cctx.Context, tx)
	})
}
