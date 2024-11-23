package httpserver

import (
	"context"
	"gateway-service/internal/authclient"
	"gateway-service/internal/config"
	currency_service "gateway-service/internal/currency-service"
	"gateway-service/internal/repository"
	"net"
	"net/http"
	"net/url"
	"time"
)

const ReadHeaderTimeoutInSeconds = 2

func NewHttpServer(ctx context.Context, cfg *config.Config) *http.Server {
	userRepository := repository.NewInMemoryUserRepository()
	userSearcher := &inMemoryUserSearcher{userRepository}
	authClientUrl := &url.URL{
		Scheme: cfg.AuthScheme,
		Host:   cfg.AuthHost + ":" + cfg.AuthPort,
	}
	authClient := authclient.NewAuthClient(userSearcher, authClientUrl)
	authHandler := authMiddleware(&authClientAuthenticator{authClient})
	currencySvcUrl := &url.URL{
		Scheme: cfg.CurrencyScheme,
		Host:   cfg.CurrencyHost + ":" + cfg.CurrencyPort,
	}
	currencySvcProxy := currency_service.NewCurrencyProxy(currencySvcUrl)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.Handle("POST /login", &loginHandler{cfg, &authTokenProvider{authClient}})
	mux.Handle("GET /exchange-rate", authHandler(&proxyHandler{&currencyRateProxyHandler{currencySvcProxy}}))
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
