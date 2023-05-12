// Package currencyapi is adapter for https://nationalbank.kz/rss/get_rates.cfm.
// Adapter works with entity.Rate.
package currencyapi

import (
	"context"
	"testing"
	"time"

	"github.com/Kin-dza-dzaa/task/internal/entity"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const testURL = "http://nationalbank.kz/rss/get_rates.cfm"

func Test_CurrrencyAPI_Rates(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Get_rates_with_valid_date",
			args: args{
				date: time.Now(),
			},
		},
		{
			name: "Get_rates_with_invalid_date",
			args: args{
				date: time.Now().Add(time.Hour * 8760),
			},
			err: entity.ErrNoRates,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		currencyAPI := New(testURL)

		t.Run(tt.name, func(t *testing.T) {
			_, err := currencyAPI.Rates(ctx, tt.args.date)
			if !cmp.Equal(err, tt.err, cmpopts.EquateErrors()) {
				t.Fatalf("want %q but got %q", tt.err, err)
			}
		})
	}
}
