package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"papertrader/internal/database"
	"papertrader/internal/helpers"

	"github.com/gorilla/websocket"
)

type TickHandler struct {
	queries *database.Queries
}

var upgrader = websocket.Upgrader{}

func NewTickHandler(queries *database.Queries) *TickHandler {
	return &TickHandler{queries: queries}
}

// TODO: Should handle various timeframes
func (h *TickHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	instrument := r.URL.Query().Get("instrument")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	start, end := calculateTickTimestamps()

	ticks, err := h.queries.GetTicks(context.Background(), database.GetTicksParams{
		Instrument: instrument,
		Start:      start,
		End:        end,
	})

	if err != nil {
		errMsg := map[string]string{
			"error":   "Failed to get ticks",
			"details": err.Error(),
		}
		if err := conn.WriteJSON(errMsg); err != nil {
			log.Println("Write error:", err)
		}
		return
	}

	for i, tick := range ticks {
		t := tick.Time
		tick.Time = time.UnixMilli(tick.Time).Truncate(time.Minute).Add(time.Minute).Unix()

		if err := conn.WriteJSON(tick); err != nil {
			log.Println("Write error:", err)
			break
		}

		if i+1 < len(ticks) {
			dt := time.Duration(ticks[i+1].Time-t) * time.Millisecond
			time.Sleep(dt)
		}
	}
}

func calculateTickTimestamps() (int64, int64) {
	now := time.Now().UTC()

	start := helpers.GetSimulatedTime(now)

	end := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)

	return start.UnixMilli(), end.UnixMilli()
}
