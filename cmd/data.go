package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"papertrader/internal/database"
	"time"

	"github.com/govalues/decimal"
)

func loadOHLCV(ctx context.Context, queries *database.Queries, instrumentID int64, instrument string) error {
	file, err := os.Open(fmt.Sprintf("tmp/@%s-Market Data Sim-GLOBEX-Futures-Minute-Trade.txt", instrument))
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	records = records[1:]

	for _, record := range records {
		parsedDate, err := time.Parse("02/01/2006", record[0])
		if err != nil {
			return err
		}

		parsedTime, err := time.Parse("15:04:05", record[1])
		if err != nil {
			return err
		}

		combined := time.Date(
			parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
			time.UTC,
		)

		open, _ := decimal.Parse(record[2])
		high, _ := decimal.Parse(record[3])
		low, _ := decimal.Parse(record[4])
		close, _ := decimal.Parse(record[5])
		volume, _ := decimal.Parse(record[6])

		err = queries.InsertOHLCV(ctx, database.InsertOHLCVParams{
			InstrumentID: instrumentID,
			Time:         combined.Unix(),
			Open:         open,
			High:         high,
			Low:          low,
			Close:        close,
			Volume:       volume,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func loadTicks(ctx context.Context, queries *database.Queries, instrumentID int64, instrument string) error {
	file, err := os.Open(fmt.Sprintf("tmp/@%s-Market Data Sim-GLOBEX-Futures-Tick-Trade.txt", instrument))
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	records = records[1:]

	for _, record := range records {
		parsedDate, err := time.Parse("02/01/2006", record[0])
		if err != nil {
			return err
		}

		parsedTime, err := time.Parse("15:04:05", record[1])
		if err != nil {
			return err
		}

		combined := time.Date(
			parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
			time.UTC,
		)

		price, _ := decimal.Parse(record[2])
		volume, _ := decimal.Parse(record[3])

		err = queries.InsertTick(ctx, database.InsertTickParams{
			InstrumentID: instrumentID,
			Time:         combined.UnixMilli(),
			Price:        price,
			Volume:       volume,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func loadInstrument(ctx context.Context, db *sql.DB, queries *database.Queries, instrument string, instrumentID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	if err := loadOHLCV(ctx, qtx, instrumentID, instrument); err != nil {
		log.Fatalf("Failed to load OHLCV data for %s: %v", instrument, err)
	}

	if err := loadTicks(ctx, qtx, instrumentID, instrument); err != nil {
		log.Fatalf("Failed to load Ticks data for %s: %v", instrument, err)
	}

	return tx.Commit()
}

func loadData(db *sql.DB, queries *database.Queries) {
	ctx := context.Background()
	instruments := []string{"NQ"}

	for _, instrument := range instruments {
		i, err := queries.GetInstrument(ctx, instrument)
		if err != nil {
			// Check the actual error value
			i, err = queries.CreateInstrument(ctx, instrument)
			if err != nil {
				log.Fatalf("Failed to load OHLCV data for %s: %v", instrument, err)
			}
		}

		loadInstrument(ctx, db, queries, instrument, i.InstrumentID)
	}
}
