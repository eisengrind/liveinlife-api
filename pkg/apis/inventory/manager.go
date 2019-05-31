package inventory

//go:generate counterfeiter -o ./mocks/manager.go . Manager
//go:generate protoc -I./../../../../../../ -I ./proto --go_out=plugins=grpc:./proto ./proto/manager.proto

import (
	"context"
	"errors"
)

// Manager is a manager for inventory objects
type Manager interface {
	Get(context.Context, Identifier) (Complete, error)
	Create(context.Context, Incomplete) (Complete, error)
	AddItem(context.Context, Identifier, *Item) error
	RemoveItem(context.Context, Identifier, *Item) error
	Delete(context.Context, Identifier) error
}

type manager struct {
	repository Repository
}

// NewManager creates a new manager for managing inventory objects
func NewManager(r Repository) Manager {
	return &manager{
		r,
	}
}

var errInvalidGUID = errors.New("error invalid guid")

func (m *manager) Get(ctx context.Context, id Identifier) (Complete, error) {
	if id.GUID() == "" {
		return nil, errInvalidGUID
	}

	return m.repository.Get(ctx, id)
}

var (
	errInvalidItemID     = errors.New("invalid item id")
	errInvalidItemAmount = errors.New("invalid item amount")
	errInvalidItemSubset = errors.New("invalid item subset")
)

func (m *manager) Create(ctx context.Context, inc Incomplete) (Complete, error) {
	for _, v := range inc.Data().Items {
		if v.ID == "" {
			return nil, errInvalidItemID
		}

		if v.Amount == 0 {
			return nil, errInvalidItemAmount
		}

		if v.Subset != -1 && !(v.Subset > 0) {
			return nil, errInvalidItemSubset
		}
	}

	c, err := m.repository.Create(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range inc.Data().Items {
		if err := m.repository.AddItem(ctx, c, v); err != nil {
			return nil, err
		}
	}

	c.Data().Items = inc.Data().Items

	return c, nil
}

func (m *manager) AddItem(ctx context.Context, id Identifier, item *Item) error {
	if id.GUID() == "" {
		return errInvalidGUID
	}

	if item.ID == "" {
		return errInvalidItemID
	}

	if item.Amount == 0 {
		return errInvalidItemAmount
	}

	if item.Subset != -1 && !(item.Subset > 0) {
		return errInvalidItemSubset
	}

	return m.repository.AddItem(ctx, id, item)
}

func (m *manager) RemoveItem(ctx context.Context, id Identifier, item *Item) error {
	if id.GUID() == "" {
		return errInvalidGUID
	}

	if item.ID == "" {
		return errInvalidItemID
	}

	if item.Amount == 0 {
		return errInvalidItemAmount
	}

	if item.Subset != -1 && !(item.Subset > 0) {
		return errInvalidItemSubset
	}

	return m.repository.RemoveItem(ctx, id, item)
}

func (m *manager) Delete(ctx context.Context, id Identifier) error {
	if id.GUID() == "" {
		return errInvalidGUID
	}

	return m.repository.Delete(ctx, id)
}
