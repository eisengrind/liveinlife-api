package serviceaccount

//go:generate counterfeiter -o ./mocks/manager.go . Manager
//go:generate protoc -I./../../../../../../ -I ./proto --go_out=plugins=grpc:./proto ./proto/manager.proto

import (
	"context"
	"errors"
	"fmt"

	"github.com/51st-state/api/pkg/rbac"
)

// Manager provides methods to manage a service account
type Manager interface {
	Repository
	GetRoles(context.Context, Identifier) (rbac.AccountRoles, error)
	SetRoles(context.Context, Identifier, rbac.AccountRoles) error
}

type manager struct {
	repository Repository
	rbac       rbac.Control
}

// NewManager creates a new service account manager
func NewManager(r Repository, rb rbac.Control) Manager {
	return &manager{r, rb}
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

func (m *manager) GetRoles(ctx context.Context, id Identifier) (rbac.AccountRoles, error) {
	if id.GUID() == "" {
		return nil, errInvalidGUID
	}

	roles, err := m.rbac.GetAccountRoles(ctx, rbac.AccountID(fmt.Sprintf(
		"service_account/%s",
		id.GUID(),
	)))
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (m *manager) SetRoles(ctx context.Context, id Identifier, roles rbac.AccountRoles) error {
	if id.GUID() == "" {
		return errInvalidGUID
	}

	return m.rbac.SetAccountRoles(ctx, rbac.AccountID(fmt.Sprintf(
		"service_account/%s",
		id.GUID(),
	)), roles)
}
