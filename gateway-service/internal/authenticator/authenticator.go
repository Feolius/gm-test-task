package authenticator

import (
	"context"
	"errors"
	"fmt"
	"gateway-service/internal/config"
	"gateway-service/internal/userstorage"
	"io"
	"net/http"
	"net/url"
)

var InvalidCredentialsError = errors.New("invalid credentials")

type UserSearcher interface {
	FindByUsernameAndPassword(ctx context.Context, username, password string) (userstorage.User, error)
}

type Authenticator struct {
	userSearcher UserSearcher
	cfg          *config.Config
}

func (a Authenticator) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.userSearcher.FindByUsernameAndPassword(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("error on attempt to load user by login credentials: %w", err)
	}
	if user.Empty() {
		return "", fmt.Errorf("%w", InvalidCredentialsError)
	}
	// @TODO put url scheme inside config as well
	authenticateRequestUrl := buildAuthUrl("http", a.cfg.AuthHost, a.cfg.AuthPort)
	req, err := http.NewRequest(http.MethodGet, authenticateRequestUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error on creating authentication request: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on requesting authentication: %w", err)
	}
	token, err := io.ReadAll(res.Body)
	fmt.Println(string(token))
	return string(token), nil
}

func NewAuthenticator(userSearcher UserSearcher, cfg *config.Config) *Authenticator {
	return &Authenticator{
		userSearcher: userSearcher,
		cfg:          cfg,
	}
}

func buildAuthUrl(proto, host, port string) string {
	u := &url.URL{
		Scheme: proto,
		Host:   host + ":" + port,
		Path:   "generate",
	}
	return u.String()
}
