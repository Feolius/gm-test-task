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

func (r *SqlCurrencyRateRepository) FindByCurrencyAndDateRange(
	ctx context.Context, currency, startDate, endDate string) ([]DailyCurrencyRate, error) {
	var res []DailyCurrencyRate
	rows, err := r.db.QueryContext(ctx, `SELECT DATE_FORMAT(day, "%Y-%m-%d") as day, currency, rate FROM currency_rates 
        WHERE currency = ? AND day >= ? AND day <= ? ORDER BY day`, currency, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r DailyCurrencyRate
		err = rows.Scan(&r.Date, &r.Currency, &r.Rate)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
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
