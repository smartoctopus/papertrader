package handlers

import (
	"net/http"

	"papertrader/internal/templates"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Page()
	err := templates.Layout(c, "Paper Trader").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
