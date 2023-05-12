// Package service implements application layer of the app.
package service

import (
	"context"
	"testing"
	"time"

	"github.com/Kin-dza-dzaa/task/internal/service/repomock"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"
)

func Test_CurrencyService_Rates(t *testing.T) {
	type args struct {
		date time.Time
		code string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		setupMock func(args args, currencyAPIMock *repomock.CurrencyAPI, currencyRepoMock *repomock.CurrencyRepo)
	}{
		{
			name: "Check_call_to_mocks",
			args: args{
				date: time.Now(),
				code: "KZT",
			},
			setupMock: func(args args, currencyAPIMock *repomock.CurrencyAPI, currencyRepoMock *repomock.CurrencyRepo) {
				currencyRepoMock.On("Rates", mock.Anything, args.date, args.code).Once().Return(nil, nil)
			},
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		currencyAPI := repomock.NewCurrencyAPI(t)
		currencyRepo := repomock.NewCurrencyRepo(t)
		currencyService := New(currencyRepo, currencyAPI, slog.Default())
		tt.setupMock(tt.args, currencyAPI, currencyRepo)

		t.Run(tt.name, func(t *testing.T) {
			_, err := currencyService.Rates(ctx, tt.args.date, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CurrencyService.Rates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
