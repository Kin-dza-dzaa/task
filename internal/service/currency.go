// Package service implements application layer of the app.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Kin-dza-dzaa/task/internal/entity"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

const defaultTaskTimeout = time.Second * 10

type (
	CurrencyRepo interface {
		Rates(ctx context.Context, date time.Time, code string) ([]entity.Rate, error)
		CreateRates(ctx context.Context, rates []entity.Rate) error
	}

	CurrencyAPI interface {
		Rates(ctx context.Context, date time.Time) ([]entity.Rate, error)
	}
)

type CurrencyService struct {
	currencyAPI      CurrencyAPI
	currencyRateRepo CurrencyRepo
	l                *slog.Logger
}

func (s *CurrencyService) Rates(ctx context.Context, date time.Time, code string) ([]entity.Rate, error) {
	return s.currencyRateRepo.Rates(ctx, date, code)
}

func (s *CurrencyService) CreateRates(ctx context.Context, date time.Time) {
	s.createRatesAsync(ctx, date)
}

func (s *CurrencyService) createRatesAsync(ctx context.Context, date time.Time) {
	// Copy req_id from http ctx.
	taskCtx, cancel := context.WithTimeout(context.Background(), defaultTaskTimeout)
	taskCtx = context.WithValue(taskCtx, middleware.RequestIDKey, middleware.GetReqID(ctx))

	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.l.ErrorCtx(
					taskCtx,
					"Panic was caught when performing a task: createRatesAsync",
					slog.Any(
						"panic val",
						err,
					),
				)
			}
		}()
		defer cancel()
		rates, err := s.currencyAPI.Rates(taskCtx, date)
		if err != nil {
			if errors.Is(err, entity.ErrNoRates) {
				return
			}
			s.l.ErrorCtx(
				taskCtx,
				"Unable to get rates from API",
				slog.String(
					"error",
					fmt.Errorf("CurrencyService - createRatesAsync - s.currencyAPI.Rate: %w", err).Error(),
				),
			)
			return
		}
		err = s.currencyRateRepo.CreateRates(taskCtx, rates)
		if err != nil {
			s.l.ErrorCtx(
				taskCtx,
				"Unable to create rate in DB",
				slog.String(
					"error",
					fmt.Errorf("CurrencyService - createRatesAsync - s.currencyRateRepo.CreateRates: %w", err).Error(),
				),
			)
		}
	}()
}

func New(currencyRepo CurrencyRepo, currencyAPI CurrencyAPI, l *slog.Logger) *CurrencyService {
	return &CurrencyService{
		l:                l,
		currencyAPI:      currencyAPI,
		currencyRateRepo: currencyRepo,
	}
}
