package handlers

import (
	"net/http"
	"papertrader/internal/templates"

	"github.com/a-h/templ"
)

func renderPage(w http.ResponseWriter, r *http.Request, contents templ.Component, title string) error {
	c := templates.Page(contents)
	err := templates.Layout(c, title).Render(r.Context(), w)

	return err
}
