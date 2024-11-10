package cron

import (
	"context"
	"currency-service/internal/config"
	"currency-service/internal/ratescollector"
	"database/sql"
	"log"

	cronlib "github.com/robfig/cron/v3"
)

func ScheduleCronJobs(ctx context.Context, cfg *config.Config, db *sql.DB) error {
	c := cronlib.New()
	_, err := c.AddFunc("@daily", func() {
		collector := ratescollector.NewRatesCollector(db)
		err := collector.Collect(ctx)
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
