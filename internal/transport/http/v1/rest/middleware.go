package rest

import (
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

func (h *CurrencyHandler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		defer func() {
			path := r.URL.Path
			h.l.InfoCtx(
				r.Context(),
				"request",
				slog.Int64("elapsed", time.Since(t).Nanoseconds()),
				slog.String("path", path),
				slog.String("method", r.Method),
			)
		}()

		next.ServeHTTP(w, r)
	})
}
