package httpserver

import (
	"context"
	"currency-service/internal/repository"
	"database/sql"
	"errors"
)

type errResponse struct {
	Message string `json:"message"`
}

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

type repositoryExchangeHistoryLoader struct {
	repository *repository.SqlCurrencyRateRepository
}

func (r *repositoryExchangeHistoryLoader) getHistory(
	ctx context.Context, currency, startDate, endDate string) ([]exchangeHistoryItem, error) {
	rates, err := r.repository.FindByCurrencyAndDateRange(ctx, currency, startDate, endDate)
	if err != nil {
		return nil, err
	}
	result := make([]exchangeHistoryItem, 0, len(rates))
	for _, rate := range rates {
		result = append(result, exchangeHistoryItem{Date: rate.Date, Rate: rate.Rate})
	}
	return result, nil
}
