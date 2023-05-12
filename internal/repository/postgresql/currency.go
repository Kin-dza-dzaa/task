// Package postgresql is adapter for postgresql.
// Adapter works with entity.Rate.
package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/Kin-dza-dzaa/task/internal/entity"
	"github.com/Kin-dza-dzaa/task/internal/service"
	"github.com/Kin-dza-dzaa/task/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

var _ = service.CurrencyRepo((*CurrencyRepo)(nil))

type CurrencyRepo struct {
	*postgres.Client
}

func (r *CurrencyRepo) Rates(ctx context.Context, date time.Time, code string) ([]entity.Rate, error) {
	query := r.Builder.Select("title", "code", "rate", "quant", "valid_at", "change").
		From("currency_rates").
		Where("valid_at = ?", date)
	if code != "" {
		query = query.Where("code = ?", code)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("CurrencyRepo - Rate - query.ToSql: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("CurrencyRepo - Rate - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var rates []entity.Rate
	for rows.Next() {
		var crntRate entity.Rate
		err := rows.Scan(
			&crntRate.Title,
			&crntRate.Code,
			&crntRate.Rate,
			&crntRate.Quant,
			&crntRate.ValidAt,
			&crntRate.Change,
		)
		if err != nil {
			return nil, fmt.Errorf("CurrencyRepo - Rate - rows.Scan: %w", err)
		}
		rates = append(rates, crntRate)
	}

	return rates, nil
}

func (r *CurrencyRepo) CreateRates(ctx context.Context, rates []entity.Rate) error {
	err := r.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, rate := range rates {
			sql, args, err := r.Builder.Insert("currency_rates").
				Columns("title", "code", "rate", "quant", "valid_at", "change").
				Values(rate.Title, rate.Code, rate.Rate, rate.Quant, rate.ValidAt, rate.Change).
				Suffix("ON CONFLICT (valid_at, code) DO NOTHING").
				ToSql()
			if err != nil {
				return fmt.Errorf("CurrencyRepo - CreateRate - r.Builder.ToSql: %w", err)
			}

			_, err = tx.Exec(ctx, sql, args...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("CurrencyRepo - CreateRate - r.Pool.BeginFunc: %w", err)
	}
	return nil
}

func New(psqlClient *postgres.Client) *CurrencyRepo {
	return &CurrencyRepo{
		Client: psqlClient,
	}
}
