package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"papertrader/internal/database"
	"papertrader/internal/helpers"
)

type OHLCVHandler struct {
	queries *database.Queries
}

func NewOHLCVHandler(queries *database.Queries) *OHLCVHandler {
	return &OHLCVHandler{queries: queries}
}

func (h *OHLCVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	instrument := r.URL.Query().Get("instrument")
	timeframe := r.URL.Query().Get("timeframe")

	w.Header().Set("Content-Type", "application/json")

	ohlcv, err := h.getOHLCV(instrument, timeframe)
	if err != nil {
		http.Error(w, "Error fetching OHLCV data", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(ohlcv)
}

func (h *OHLCVHandler) getOHLCV(instrument string, timeframe string) ([]database.GetOHLCVRow, error) {
	allowed_instruments := map[string]struct{}{
		"NQ": {},
		"ES": {},
		"GC": {},
	}

	if _, exists := allowed_instruments[instrument]; !exists {
		return nil, errors.New("Instrument not allowed")
	}

	duration, err := helpers.TimeframeToDuration(timeframe)
	if err != nil {
		return nil, err
	}

	start, end, ticks_start, ticks_end := calculateOHLCVTimestamps()

	ohlcv, err := h.queries.GetOHLCV(context.Background(), database.GetOHLCVParams{
		Instrument: instrument,
		Start:      start,
		End:        end,
	})
	if err != nil {
		return nil, err
	}

	ticks, err := h.queries.GetTicks(context.Background(), database.GetTicksParams{
		Instrument: instrument,
		Start:      ticks_start,
		End:        ticks_end,
	})
	if err != nil {
		return nil, err
	}

	if len(ticks) != 0 {
		tick := ticks[0]

		lastCandle := database.GetOHLCVRow{
			Time:   time.UnixMilli(tick.Time).Truncate(time.Minute).Add(time.Minute).Unix(),
			Open:   tick.Price,
			High:   tick.Price,
			Low:    tick.Price,
			Close:  tick.Price,
			Volume: tick.Volume,
		}

		for i := 1; i < len(ticks); i++ {
			tick := ticks[i]

			lastCandle.Close = tick.Price
			lastCandle.Volume, _ = lastCandle.Volume.Add(tick.Volume)

			if tick.Price.Cmp(lastCandle.Low) < 0 {
				lastCandle.Low = tick.Price
			}

			if tick.Price.Cmp(lastCandle.High) > 0 {
				lastCandle.High = tick.Price
			}
		}

		ohlcv = append(ohlcv, lastCandle)
	}

	return recalculateOHLCVOnTimeframe(ohlcv, duration), nil
}

func calculateOHLCVTimestamps() (int64, int64, int64, int64) {
	now := time.Now().UTC()
	startOfToday := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)

	numberOfDays := time.Duration(5)
	start := startOfToday.Add(-time.Hour * 24 * (numberOfDays + 1)) // Need to add one day, because the simulation runs on yesterday data

	simulated_time := helpers.GetSimulatedTime(now)
	end := simulated_time.Truncate(time.Minute)

	return start.Unix(), end.Unix(), end.UnixMilli(), simulated_time.UnixMilli()
}

func createCandle(candle database.GetOHLCVRow, timeframe time.Duration) database.GetOHLCVRow {
	t := time.Unix(candle.Time, 0)

	return database.GetOHLCVRow{
		Time:   helpers.CandleTime(t, timeframe).Unix(),
		Open:   candle.Open,
		High:   candle.High,
		Low:    candle.Low,
		Close:  candle.Close,
		Volume: candle.Volume,
	}
}

func updateCandle(oldCandle database.GetOHLCVRow, newCandle database.GetOHLCVRow) database.GetOHLCVRow {
	if newCandle.High.Cmp(oldCandle.High) > 0 {
		oldCandle.High = newCandle.High
	}

	if newCandle.Low.Cmp(oldCandle.Low) < 0 {
		oldCandle.Low = newCandle.Low
	}

	oldCandle.Close = newCandle.Close
	oldCandle.Volume, _ = oldCandle.Volume.Add(newCandle.Volume)

	return oldCandle
}

func recalculateOHLCVOnTimeframe(ohlcv []database.GetOHLCVRow, timeframe time.Duration) []database.GetOHLCVRow {
	result := []database.GetOHLCVRow{}

	result = append(result, createCandle(ohlcv[0], timeframe))

	ohlcv = ohlcv[1:]
	for _, candle := range ohlcv {
		t := time.Unix(candle.Time, 0)

		if helpers.CandleTime(t, timeframe).Unix() == result[len(result)-1].Time {
			result[len(result)-1] = updateCandle(result[len(result)-1], candle)
		} else {
			result = append(result, createCandle(candle, timeframe))
		}
	}

	return result
}
