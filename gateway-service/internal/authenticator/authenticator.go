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
	"strings"
)

const generateTokenPath = "generate"
const validateTokenPath = "validate"

var InvalidCredentialsError = errors.New("invalid credentials")
var AuthServiceTokenGenError = errors.New("auth service token generation failed")
var TokenNotProvidedError = errors.New("token not provided")
var InvalidTokenFormatError = errors.New("invalid token format")
var AuthServiceTokenNotProvidedError = errors.New("auth token not provided for authentication service")
var AuthServiceTokenExpired = errors.New("auth service token expired")
var AuthServiceUnexpectedError = errors.New("auth service unexpected error")

type UserSearcher interface {
	FindByUsernameAndPassword(ctx context.Context, username, password string) (userstorage.User, error)
}

type Authenticator struct {
	userSearcher UserSearcher
	cfg          *config.Config
}

func (a *Authenticator) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.userSearcher.FindByUsernameAndPassword(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("error on attempt to load user by login credentials: %w", err)
	}
	if user.Empty() {
		return "", fmt.Errorf("%w", InvalidCredentialsError)
	}
	// @TODO put url scheme inside config as well
	genTokenUrl := buildAuthServiceUrl("http", a.cfg.AuthHost, a.cfg.AuthPort, generateTokenPath)
	req, err := http.NewRequest(http.MethodGet, genTokenUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error on creating generate token request: %w", err)
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on generate token request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w", AuthServiceTokenGenError)
	}
	token, err := io.ReadAll(res.Body)
	return string(token), nil
}

// Authenticate Ideally this method must return User.
func (a *Authenticator) Authenticate(ctx context.Context, req *http.Request) error {
	token := req.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("%w", TokenNotProvidedError)
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return fmt.Errorf("%w", InvalidTokenFormatError)
	}
	token = strings.TrimPrefix(token, "Bearer ")
	validateTokenUrl := buildAuthServiceUrl("http", a.cfg.AuthHost, a.cfg.AuthPort, validateTokenPath)
	req, err := http.NewRequest(http.MethodGet, validateTokenUrl, nil)
	if err != nil {
		return fmt.Errorf("error on creating authentication request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on authentication request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusBadRequest {
			return fmt.Errorf("%w", AuthServiceTokenNotProvidedError)
		}
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("%w", AuthServiceTokenExpired)
		}
		return fmt.Errorf("%w", AuthServiceUnexpectedError)
	}
	return nil
}

func NewAuthenticator(userSearcher UserSearcher, cfg *config.Config) *Authenticator {
	return &Authenticator{
		userSearcher: userSearcher,
		cfg:          cfg,
	}
}

func buildAuthServiceUrl(proto, host, port, path string) string {
	u := &url.URL{
		Scheme: proto,
		Host:   host + ":" + port,
		Path:   path,
	}
	return u.String()
}
