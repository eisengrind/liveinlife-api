package serviceaccount

//go:generate counterfeiter -o ./mocks/manager.go . Manager
//go:generate protoc -I ./proto --go_out=plugins=grpc:./proto ./proto/manager.proto

import (
	"context"
	"errors"
)

// Manager provides methods to manage a service account
type Manager interface {
	Repository
}

type manager struct {
	repository Repository
}

// NewManager creates a new service account manager
func NewManager(r Repository) Manager {
	return &manager{r}
}

var errInvalidGUID = errors.New("invalid guid given")

func (m *manager) Get(ctx context.Context, id Identifier) (Complete, error) {
	if id.GUID() == "" {
		return nil, errInvalidGUID
	}

	return m.repository.Get(ctx, id)
}

var errInvalidName = errors.New("invalid name given")

func (m *manager) Update(ctx context.Context, c Complete) error {
	if c.GUID() == "" {
		return errInvalidGUID
	}

	if c.Data().Name == "" {
		return errInvalidName
	}

	return m.repository.Update(ctx, c)
}

func (m *manager) Create(ctx context.Context, inc Incomplete) (Complete, error) {
	if inc.Data().Name == "" {
		return nil, errInvalidName
	}

	return m.repository.Create(ctx, inc)
}

func (m *manager) Delete(ctx context.Context, id Identifier) error {
	if id.GUID() == "" {
		return errInvalidGUID
	}

	return m.repository.Delete(ctx, id)
}
