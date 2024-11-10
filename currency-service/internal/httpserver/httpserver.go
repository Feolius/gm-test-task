package httpserver

import (
	"context"
	"currency-service/internal/config"
	"database/sql"
	"net"
	"net/http"
	"time"
)

const ReadHeaderTimeoutInSeconds = 2

type ExchangeRateHandler struct{}

func (h *ExchangeRateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello"))
}

func NewHttpServer(ctx context.Context, cfg *config.Config, dbConn *sql.DB) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.Handle("GET /exchange-rate", &ExchangeRateHandler{})
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
