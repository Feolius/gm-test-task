package httpserver

import (
	"context"
	"errors"
	"fmt"
	"gateway-service/internal/authclient"
	currency_service "gateway-service/internal/currency-service"
	"gateway-service/internal/repository"
	"net/http"
	"strings"
)

var ErrTokenNotProvided = errors.New("token not provided")
var ErrInvalidTokenFormat = errors.New("invalid token format")

type invalidCredentialsError struct {
	username string
	err      error
}

func (e *invalidCredentialsError) Error() string {
	return fmt.Sprintf("invalid auth credentials for user %s", e.username)
}

func (e *invalidCredentialsError) Unwrap() error {
	return e.err
}

type authTokenProvider struct {
	*authclient.AuthClient
}

func (p *authTokenProvider) getToken(ctx context.Context, username, password string) (string, error) {
	token, err := p.AuthClient.GetToken(ctx, username, password)
	if err != nil {
		if errors.Is(err, authclient.ErrInvalidCredentials) {
			return "", &invalidCredentialsError{username, err}
		}
		return "", fmt.Errorf("failed to get token for user %s: %w", username, err)
	}
	return token, nil
}

type inMemoryUserSearcher struct {
	repository *repository.InMemoryUserRepository
}

func (s *inMemoryUserSearcher) FindByUsernameAndPassword(
	ctx context.Context, username, password string,
) (authclient.User, error) {
	storageUser, err := s.repository.FindByUsernameAndPassword(ctx, username, password)
	if err != nil {
		return authclient.User{}, err
	}
	return authclient.User{Id: storageUser.Id}, nil
}

type authClientAuthenticator struct {
	*authclient.AuthClient
}

func (a *authClientAuthenticator) authenticate(req *http.Request) error {
	token := req.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("%w", ErrTokenNotProvided)
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return fmt.Errorf("%w", ErrInvalidTokenFormat)
	}
	token = strings.TrimPrefix(token, "Bearer ")
	err := a.AuthClient.CheckToken(req.Context(), token)
	if err != nil {
		return err
	}
	return nil
}

type requestProxy interface {
	handle(writer http.ResponseWriter, request *http.Request) error
}

type currencyRateProxyHandler struct {
	currencySvcProxy *currency_service.CurrencyServiceProxy
}

func (h *currencyRateProxyHandler) handle(writer http.ResponseWriter, request *http.Request) error {
	return h.currencySvcProxy.ExchangeRate(writer, request)
}

type historyRateProxyHandler struct {
	currencySvcProxy *currency_service.CurrencyServiceProxy
}

func (h *historyRateProxyHandler) handle(writer http.ResponseWriter, request *http.Request) error {
	return h.currencySvcProxy.ExchangeHistory(writer, request)
}
