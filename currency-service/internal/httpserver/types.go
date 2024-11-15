package httpserver

import (
	"context"
	"currency-service/internal/repository"
	"database/sql"
	"errors"
)

type repositoryCurrencyRateLoader struct {
	repository *repository.SqlCurrencyRateRepository
}

func (r *repositoryCurrencyRateLoader) load(ctx context.Context, date, currency string) (float64, error) {
	rate, err := r.repository.Find(ctx, date, currency)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	if rate == nil {
		return 0, nil
	}
	return rate.Rate, nil
}
