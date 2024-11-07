package authclient

import (
	"context"
	"errors"
	"fmt"
	"gateway-service/internal/config"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const generateTokenPath = "generate"
const validateTokenPath = "validate"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrAuthServiceTokenGen = errors.New("auth service token generation failed")
var ErrInvalidTokenFormat = errors.New("invalid token format")
var ErrAuthServiceTokenNotProvided = errors.New("auth token not provided for authentication service")
var ErrAuthServiceTokenExpired = errors.New("auth service token expired")
var ErrAuthServiceUnexpectedError = errors.New("auth service unexpected error")

type User struct {
	Id int
}

func (u *User) empty() bool {
	return u.Id == 0
}

type UserSearcher interface {
	FindByUsernameAndPassword(ctx context.Context, username, password string) (User, error)
}

type AuthClient struct {
	userSearcher UserSearcher
	cfg          *config.Config
}

func (a *AuthClient) GetToken(ctx context.Context, username, password string) (string, error) {
	user, err := a.userSearcher.FindByUsernameAndPassword(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("error on attempt to load user by login credentials: %w", err)
	}
	if user.empty() {
		return "", fmt.Errorf("%w", ErrInvalidCredentials)
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
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w", ErrAuthServiceTokenGen)
	}
	token, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error on read token response body: %w", err)
	}
	return string(token), nil
}

// CheckToken Ideally should return User.
func (a *AuthClient) CheckToken(ctx context.Context, token string) error {
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
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusBadRequest {
			return fmt.Errorf("%w", ErrAuthServiceTokenNotProvided)
		}
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("%w", ErrAuthServiceTokenExpired)
		}
		return fmt.Errorf("%w", ErrAuthServiceUnexpectedError)
	}
	return nil
}

func NewAuthClient(userSearcher UserSearcher, cfg *config.Config) *AuthClient {
	return &AuthClient{
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
