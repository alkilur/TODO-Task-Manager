package http_server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func (c *Controller) GetNextDate(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))
		// TODO: log something?

		now, err := time.Parse(TimeLayout, r.FormValue("now"))
		if err != nil {
			c.sendError(w, r, fmt.Errorf("invalid 'now' format"))
			return
		}

		nextDate, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
		if err != nil {
			c.sendError(w, r, err)
			return
		}

		render.PlainText(w, r, nextDate)
	}
}
