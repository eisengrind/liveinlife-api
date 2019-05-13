package auth

import (
	"context"
	"time"
)

// Repository for the authentication service
//go:generate counterfeiter -o ./mocks/repository.go . Repository
type Repository interface {
	LoginAttemptsCountSince(ctx context.Context, id string, t time.Time) (uint64, error)
	AddLoginAttempt(ctx context.Context, id string, t time.Time) error
}
