package httpserver

import (
	"context"
	"currency-service/internal/config"
	"currency-service/internal/repository"
	"database/sql"
	"net"
	"net/http"
	"time"
)

const ReadHeaderTimeoutInSeconds = 2

func NewHttpServer(ctx context.Context, cfg *config.Config, db *sql.DB) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	currencyRateRepository := repository.NewSqlCurrencyRateRepository(db)
	mux.Handle("GET /exchange-rate", &exchangeRateHandler{&repositoryCurrencyRateLoader{currencyRateRepository}})
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: ReadHeaderTimeoutInSeconds * time.Second,
	}

	return httpServer
}
