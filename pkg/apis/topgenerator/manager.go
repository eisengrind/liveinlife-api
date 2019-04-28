package top

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/51st-state/api/pkg/problems"
	"github.com/pkg/errors"
)

// Manager for tops
type Manager struct {
	repository Repository
}

// NewManager for top objects
func NewManager(r Repository) *Manager {
	return &Manager{
		repository: r,
	}
}

var (
	errInvalidSex = errors.New("invalid sex given")
)

var errTopNotFound = problems.New("top not found", "the given ids are not linked to a top", http.StatusNotFound)

// Get top information
func (m *Manager) Get(ctx context.Context, id Identifier) (Complete, error) {
	if id.Sex() > 1 {
		return nil, errInvalidSex
	}

	c, err := m.repository.Get(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errTopNotFound
		}

		return nil, err
	}

	return c, nil
}

// Upsert a top into a repository
func (m *Manager) Upsert(ctx context.Context, c Complete) error {
	if c.Sex() > 1 {
		return errInvalidSex
	}

	return m.repository.Upsert(ctx, c)
}
