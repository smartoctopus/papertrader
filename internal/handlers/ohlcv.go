package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"papertrader/internal/database"

	"github.com/govalues/decimal"
)

type OHLCV struct {
	Time   int64           `json:"time"`
	Open   decimal.Decimal `json:"open"`
	High   decimal.Decimal `json:"high"`
	Low    decimal.Decimal `json:"low"`
	Close  decimal.Decimal `json:"close"`
	Volume decimal.Decimal `json:"volume"`
}

type OHLCVHandler struct {
	queries *database.Queries
}

func NewOHLCVHandler(queries *database.Queries) *OHLCVHandler {
	return &OHLCVHandler{queries: queries}
}

func (h *OHLCVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	instrument := r.URL.Query().Get("instrument")

	w.Header().Set("Content-Type", "application/json")

	// TODO: Make this function a DB request
	ohlcv, err := generateOHLCV(instrument)
	if err != nil {
		http.Error(w, "Error fetching OHLCV data", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(ohlcv)
}

// TODO: Load this data to a DB
func generateOHLCV(instrument string) ([]OHLCV, error) {
	allowed_instruments := map[string]struct{}{
		"NQ": {},
		"ES": {},
		"GC": {},
	}

	if _, exists := allowed_instruments[instrument]; !exists {
		return nil, errors.New("Instrument not allowed")
	}

	file, err := os.Open(fmt.Sprintf("tmp/%s.csv", instrument))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := []OHLCV{}

	for _, record := range records {
		parsedDate, err := time.Parse("02/01/2006", record[0])
		if err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse("15:04:05", record[1])
		if err != nil {
			return nil, err
		}

		combined := time.Date(
			parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
			parsedDate.Location(),
		)

		open, _ := decimal.Parse(record[2])
		high, _ := decimal.Parse(record[3])
		low, _ := decimal.Parse(record[4])
		close, _ := decimal.Parse(record[5])
		volume, _ := decimal.Parse(record[6])

		result = append(result, OHLCV{
			Time:   combined.Unix(),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	return result, nil
}
