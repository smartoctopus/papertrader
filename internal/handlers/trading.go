package handlers

import (
	"net/http"
	"papertrader/internal/database"
	"papertrader/internal/templates"
)

type TradingHandler struct {
	queries *database.Queries
}

func NewTradingHandler(queries *database.Queries) *TradingHandler {
	return &TradingHandler{queries: queries}
}

func (h *TradingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := templates.Trading().Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
