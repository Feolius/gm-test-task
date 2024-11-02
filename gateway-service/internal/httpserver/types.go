package httpserver

import (
	"context"
	"errors"
	"fmt"
	"gateway-service/internal/authenticator"
	"gateway-service/internal/userstorage"
)

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
	*authenticator.Authenticator
}

func (p *authTokenProvider) getToken(ctx context.Context, username, password string) (string, error) {
	token, err := p.Authenticator.GetToken(ctx, username, password)
	if err != nil {
		if errors.Is(err, authenticator.ErrInvalidCredentials) {
			return "", &invalidCredentialsError{username, err}
		}
		return "", fmt.Errorf("failed to get token for user %s: %w", username, err)
	}
	return token, nil
}

type inMemoryUserSearcher struct {
	storage *userstorage.InMemoryStorage
}

func (s *inMemoryUserSearcher) FindByUsernameAndPassword(
	ctx context.Context, username, password string,
) (userstorage.User, error) {
	return s.storage.FindByUsernameAndPassword(ctx, username, password)
}
