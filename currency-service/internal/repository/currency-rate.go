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

func (r *SqlCurrencyRateRepository) Find(ctx context.Context, date, currency string) (*DailyCurrencyRate, error) {
	var res DailyCurrencyRate
	err := r.db.QueryRowContext(ctx, "SELECT day, currency, rate FROM currency_rates WHERE day = ? and currency = ?",
		date, currency).Scan(&res.Date, &res.Currency, &res.Rate)
	if err != nil {
		return nil, err
	}
	return &res, nil
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
