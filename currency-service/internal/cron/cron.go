package cron

import (
	"context"
	"currency-service/internal/config"
	"currency-service/internal/ratescollector"
	"currency-service/internal/repository"
	"database/sql"
	"log"
	"time"

	cronlib "github.com/robfig/cron/v3"
)

const CollectorTimeoutInSeconds = 30

type ratesCollectorSaver struct {
	repository *repository.SqlCurrencyRateRepository
}

func (s *ratesCollectorSaver) Save(ctx context.Context, cr *ratescollector.CurrencyRate) error {
	dcr := &repository.DailyCurrencyRate{
		Date:     cr.Date,
		Currency: cr.Currency,
		Rate:     cr.Rate,
	}
	return s.repository.Save(ctx, dcr)
}

func ScheduleCronJobs(ctx context.Context, cfg *config.Config, db *sql.DB) error {
	c := cronlib.New()
	repo := repository.NewSqlCurrencyRateRepository(db)
	collector := ratescollector.NewRatesCollector(&ratesCollectorSaver{repo})
	_, err := c.AddFunc("@daily", func() {
		timeoutCtx, cancel := context.WithTimeout(ctx, CollectorTimeoutInSeconds*time.Second)
		defer cancel()
		err := collector.Collect(timeoutCtx)
		if err != nil {
			log.Printf("error on attempt to collect currency rates: %v", err)
		}
	})
	if err != nil {
		return err
	}
	c.Start()
	go func() {
		<-ctx.Done()
		c.Stop()
	}()
	return nil
}
