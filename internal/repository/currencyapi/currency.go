// Package currencyapi is adapter for https://nationalbank.kz/rss/get_rates.cfm.
// Adapter works with entity.Rate.
package currencyapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Kin-dza-dzaa/task/internal/entity"
	"github.com/Kin-dza-dzaa/task/internal/service"
	"github.com/shopspring/decimal"
)

type Rates struct {
	Items []struct {
		Quant   uint            `xml:"quant"`
		Title   string          `xml:"fullname"`
		Code    string          `xml:"title"`
		Change  string          `xml:"change"`
		Rate    decimal.Decimal `xml:"description"`
		ValidAt time.Time
	} `xml:"item"`
}

var _ = service.CurrencyAPI((*CurrrencyAPI)(nil))

const (
	dateQueryKey = "fdate"
	dateFormat   = "02.01.2006"
)

type CurrrencyAPI struct {
	targetURL string
	client    *http.Client
}

func (r *CurrrencyAPI) Rates(ctx context.Context, date time.Time) ([]entity.Rate, error) {
	u, err := url.Parse(r.targetURL)
	if err != nil {
		return nil, fmt.Errorf("CurrrencyAPI - Rate - url.Parse: %w", err)
	}
	q := u.Query()
	q.Add(dateQueryKey, date.Format(dateFormat))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("CurrrencyAPI - Rate - http.NewRequest: %w", err)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CurrrencyAPI - Rate - r.client.Do: %w", err)
	}
	defer res.Body.Close()

	var rates Rates
	if err := xml.NewDecoder(res.Body).Decode(&rates); err != nil {
		return nil, fmt.Errorf("CurrrencyAPI - Rates - xml.NewDecoder.Decode: %w", err)
	}
	if len(rates.Items) == 0 {
		return nil, entity.ErrNoRates
	}

	ratesWithDate := make([]entity.Rate, 0, len(rates.Items))
	for _, v := range rates.Items {
		v.ValidAt = date
		ratesWithDate = append(ratesWithDate, entity.Rate(v))
	}

	return ratesWithDate, nil
}

func New(addr string) *CurrrencyAPI {
	return &CurrrencyAPI{
		targetURL: addr,
		client:    &http.Client{},
	}
}
