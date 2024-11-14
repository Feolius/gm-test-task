package repository

import (
	"context"
	"database/sql"
)

type DailyCurrencyRate struct {
	Date     string
	Currency string
	Rate     float64
}

type SqlCurrencyRateRepository struct {
	db *sql.DB
}

func (r *SqlCurrencyRateRepository) Save(ctx context.Context, rate *DailyCurrencyRate) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO currency_rates (day, currency, rate) VALUES (?, ?, ?) 
					ON DUPLICATE KEY UPDATE rate = ?`, rate.Date, rate.Currency, rate.Rate, rate.Rate)
	if err != nil {
		return err
	}
	return nil
}

func NewSqlCurrencyRateRepository(db *sql.DB) *SqlCurrencyRateRepository {
	return &SqlCurrencyRateRepository{db}
}
