package key

//go:generate protoc -I ./proto --go_out=plugins=grpc:./proto ./proto/manager.proto

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"

	"github.com/51st-state/api/pkg/apis/serviceaccount"
)

// Manager for service account keys
type Manager interface {
	Get(context.Context, Identifier) (Complete, error)
	Update(context.Context, Complete) error
	Delete(context.Context, Identifier) error

	Create(context.Context, Incomplete) (*ClientKey, error)
}

type manager struct {
	repository     Repository
	serviceaccount serviceaccount.Manager
}

// NewManager creates a new service account key manager
func NewManager(r Repository, s serviceaccount.Manager) Manager {
	return &manager{r, s}
}

var errInvalidGUID = errors.New("invalid guid")

func (m *manager) Get(ctx context.Context, id Identifier) (Complete, error) {
	if id.GUID() == "" {
		return nil, errInvalidGUID
	}

	return m.repository.Get(ctx, id)
}

var errInvalidName = errors.New("invalid key name")

func (m *manager) Update(ctx context.Context, c Complete) error {
	if c.GUID() == "" {
		return errInvalidGUID
	}

	if c.Data().Name == "" {
		return errInvalidName
	}

	return m.repository.Update(ctx, c)
}

func (m *manager) Delete(ctx context.Context, id Identifier) error {
	if id.GUID() == "" {
		return errInvalidGUID
	}

	return m.repository.Delete(ctx, id)
}

const privateKeyBitSize = 2048

func (m *manager) Create(ctx context.Context, inc Incomplete) (*ClientKey, error) {
	if _, err := m.serviceaccount.Get(
		ctx,
		serviceaccount.NewIdentifier(inc.Data().ServiceAccountGUID),
	); err != nil {
		return nil, err
	}

	if inc.Data().Name == "" {
		return nil, errInvalidName
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeyBitSize)
	if err != nil {
		return nil, err
	}

	inc.Data().PublicKey = &privateKey.PublicKey

	c, err := m.repository.Create(ctx, inc)
	if err != nil {
		return nil, err
	}

	return &ClientKey{
		ServiceAccountGUID: c.Data().ServiceAccountGUID,
		GUID:               c.GUID(),
		PrivateKey:         privateKey,
	}, nil
}
