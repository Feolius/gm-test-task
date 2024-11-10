package main

import (
	"context"
	migrations "currency-service/db"
	"currency-service/internal/config"
	"currency-service/internal/cron"
	"currency-service/internal/httpserver"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.GetConfig()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	err := migrations.RunMigrations(ctx, dbstring)
	if err != nil {
		log.Panicf("error running migrations: %v", err)
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Panicf("error connecting to database: %v", err)
	}
	err = cron.ScheduleCronJobs(ctx, cfg, db)
	if err != nil {
		log.Panicf("error scheduling cron jobs: %v", err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	httpServer := httpserver.NewHttpServer(ctx, cfg, db)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err = g.Wait(); err != nil {
		log.Printf("exit reason: %s \n", err)
	}
}
