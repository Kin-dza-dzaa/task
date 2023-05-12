// Package rest is port layer of the app for currency.
package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Kin-dza-dzaa/task/config"
	"github.com/Kin-dza-dzaa/task/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"golang.org/x/exp/slog"
)

const (
	dateFormat  = "02.01.2006"
	datePathKey = "date"
	codePathKey = "code"
)

type GetRatesResponse struct {
	Rates []entity.Rate `json:"rates"`
}

type CurrencyService interface {
	Rates(ctx context.Context, date time.Time, code string) ([]entity.Rate, error)
	CreateRates(ctx context.Context, date time.Time)
}

//	@title			Rate API
//	@version		0.1
//	@description	REST API for currency rates to KZT.

//	@contact.name	API Support

//	@host		localhost:8000
//	@BasePath	/v1
type CurrencyHandler struct {
	currencyService CurrencyService
	l               *slog.Logger
}

// List rates by date and code
//
//	@Summary	Gets all rates by given date and code.
//	@Tags		rates
//	@Produce	json
//	@Success	200		{object}	GetRatesResponse	"Rates"
//
//	@Param		date	path		string				true	"Date in format: 11.11.2011"
//	@Param		code	path		string				true	"Code of currency: KZT"
//
//	@Failure	400		{object}	httpResponse		"Invalid date"
//	@Failure	500		{object}	httpResponse		"Internal error"
//	@Router		/rates/{date}/{code} [get]
func (h *CurrencyHandler) Rates(w http.ResponseWriter, r *http.Request) {
	pathDate := chi.URLParam(r, datePathKey)
	code := chi.URLParam(r, codePathKey)
	date, err := time.Parse(dateFormat, pathDate)
	if err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Message: http.StatusText(http.StatusBadRequest),
				Path:    r.URL.Path,
			},
		)
		return
	}

	rates, err := h.currencyService.Rates(r.Context(), date, code)
	if err != nil {
		h.l.ErrorCtx(
			r.Context(),
			"Unable to get rates from service",
			slog.String(
				"error",
				fmt.Errorf("CurrencyHandler - Rates - h.currencyService.Rates: %w", err).Error(),
			),
		)
		h.encode(
			w,
			http.StatusInternalServerError,
			httpResponse{
				Message: http.StatusText(http.StatusInternalServerError),
				Path:    r.URL.Path,
			},
		)
		return
	}

	h.encode(
		w,
		http.StatusOK,
		GetRatesResponse{
			Rates: rates,
		},
	)
}

// List rates by date
//
//	@Summary	Gets all rates by given date.
//	@Tags		rates
//	@Produce	json
//	@Success	200		{object}	GetRatesResponse	"Rates"
//
//	@Param		date	path		string				true	"Date in format: 11.11.2011"
//
//	@Failure	400		{object}	httpResponse		"Invalid date"
//	@Failure	500		{object}	httpResponse		"Internal error"
//	@Router		/rates/{date} [get]
func RatesByDate() {}

// Create rates
//
//	@Summary	Makes call to external API and populates DB asynchronously.
//	@Tags		rates
//	@Produce	json
//	@Success	202		{object}	httpResponse	"Request is accepted and is being proccessed"
//
//	@Param		date	path		string			true	"Date in format: 11.11.2011"
//
//	@Failure	400		{object}	httpResponse	"Invalid date"
//	@Router		/rates/{date} [post]
func (h *CurrencyHandler) CreateRates(w http.ResponseWriter, r *http.Request) {
	pathDate := chi.URLParam(r, datePathKey)
	date, err := time.Parse(dateFormat, pathDate)
	if err != nil {
		h.encode(
			w,
			http.StatusBadRequest,
			httpResponse{
				Message: http.StatusText(http.StatusBadRequest),
				Path:    r.URL.Path,
			},
		)
		return
	}

	h.currencyService.CreateRates(r.Context(), date)

	h.encode(
		w,
		http.StatusAccepted,
		httpResponse{
			Message: http.StatusText(http.StatusAccepted),
			Path:    r.URL.Path,
		},
	)
}

func (h *CurrencyHandler) Register(c *chi.Mux, cfg config.Config) {
	c.Use(middleware.RequestID)
	c.Use(h.logRequest)
	c.Use(middleware.Recoverer)
	c.Use(httprate.Limit(
		cfg.HTTPServer.RateLimit,
		cfg.HTTPServer.RateWindow,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}),
	))

	c.Get("/swagger/*", httpSwagger.Handler())

	c.Route("/v1", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Post("/rates/{date:\\d\\d\\.\\d\\d\\.\\d\\d\\d\\d}", h.CreateRates)
		r.Get("/rates/{date:\\d\\d\\.\\d\\d\\.\\d\\d\\d\\d}/{code}", h.Rates)
		r.Get("/rates/{date:\\d\\d\\.\\d\\d\\.\\d\\d\\d\\d}", h.Rates)
	})
}

func New(l *slog.Logger, currencyService CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
		l:               l,
	}
}
