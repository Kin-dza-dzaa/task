package postgresql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Kin-dza-dzaa/task/internal/entity"
	"github.com/Kin-dza-dzaa/task/pkg/logger"
	"github.com/Kin-dza-dzaa/task/pkg/postgres"
	"github.com/adrianbrad/psqldocker"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v4"
	"golang.org/x/exp/slog"
)

const SQLUP = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS currency_rates (
		id          UUID            DEFAULT uuid_generate_v4(),
		title       TEXT                                        NOT NULL,
		code        TEXT                                        NOT NULL,
		rate        NUMERIC(18, 2)                              NOT NULL,
		quant       INTEGER         DEFAULT 1,
		change      TEXT                                        NOT NULL, 
		valid_at    DATE                                        NOT NULL,
		PRIMARY KEY(id),
		UNIQUE (valid_at, code)
);
`

func Test_CurrencyRepo_Create_Rates(t *testing.T) {
	type args struct {
		rates []entity.Rate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create_some_rates",
			args: args{
				rates: []entity.Rate{
					{
						Quant:  123,
						Title:  "some_title",
						Code:   "123",
						Change: "123",
					},
					{
						Quant:  123,
						Title:  "some_title",
						Code:   "123",
						Change: "123",
					},
					{
						Quant:  1234,
						Title:  "some_title_1",
						Code:   "1234",
						Change: "1234",
					},
					{
						Quant:  1235,
						Title:  "some_title_2",
						Code:   "1235",
						Change: "1235",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		currencyRepo := setupCurrencyRepo(ctx, t, tt.name)

		t.Run(tt.name, func(t *testing.T) {
			err := currencyRepo.CreateRates(ctx, tt.args.rates)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err: %v but got: %v", tt.wantErr, err)
			}
		})
	}
}

func Test_CurrencyRepo_Rates(t *testing.T) {
	type args struct {
		date time.Time
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    []entity.Rate
		wantErr bool
	}{
		{
			name: "Get_empty_rates",
		},
		{
			name: "Get_needed_rates",
			want: []entity.Rate{
				{
					Quant:  123,
					Title:  "some_title",
					Code:   "123",
					Change: "123",
				},
				{
					Quant:  1234,
					Title:  "some_title_1",
					Code:   "1234",
					Change: "1234",
				},
				{
					Quant:  1235,
					Title:  "some_title_2",
					Code:   "1235",
					Change: "1235",
				},
			},
		},
		{
			name: "Get_needed_rates_with_code",
			args: args{
				code: "123",
			},
			want: []entity.Rate{
				{
					Quant:  123,
					Title:  "some_title",
					Code:   "123",
					Change: "123",
				},
			},
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		currencyRepo := setupCurrencyRepoRates(ctx, t, tt.name, tt.want)

		t.Run(tt.name, func(t *testing.T) {
			got, err := currencyRepo.Rates(ctx, tt.args.date, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Fatalf("want err: %v but got: %v", tt.wantErr, err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("want: %v but got: %v diff: %v", tt.want, got, diff)
			}
		})
	}
}

func setupCurrencyRepoRates(ctx context.Context, t *testing.T, containerName string, rates []entity.Rate) *CurrencyRepo {
	t.Helper()
	currencyRepo := setupCurrencyRepo(ctx, t, containerName)

	err := currencyRepo.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, rate := range rates {
			sql, args, err := currencyRepo.Builder.Insert("currency_rates").
				Columns("title", "code", "rate", "quant", "valid_at", "change").
				Values(rate.Title, rate.Code, rate.Rate, rate.Quant, rate.ValidAt, rate.Change).
				Suffix("ON CONFLICT (valid_at, code) DO NOTHING").
				ToSql()
			if err != nil {
				return fmt.Errorf("currencyRepo.Builder.ToSql: %w", err)
			}

			_, err = tx.Exec(ctx, sql, args...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("setupCurrencyRepoRates - currencyRepo.Pool.BeginFunc: %v", err)
	}

	return currencyRepo
}

func setupCurrencyRepo(ctx context.Context, t *testing.T, containerName string) *CurrencyRepo {
	t.Helper()
	c, err := psqldocker.NewContainer(
		"testuser",
		"12345",
		"test",
		psqldocker.WithSQL(SQLUP),
		psqldocker.WithContainerName(containerName),
	)
	if err != nil {
		t.Fatalf("setupPostgresql: %v", err)
	}
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Fatalf("setupPostgresql: %v", err)
		}
	})

	l := logger.NewPGXLogger(slog.Default())
	poolClient, err := postgres.NewClient(
		ctx,
		fmt.Sprintf("postgresql://testuser:12345@localhost:%v/test", c.Port()),
		10,
		pgx.LogLevelDebug,
		l,
	)
	if err != nil {
		t.Fatalf("setupPostgresql - postgres.NewClient: %v", err)
	}
	return New(poolClient)
}
