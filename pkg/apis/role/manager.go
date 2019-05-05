package role

import (
	"context"
	"errors"
	"regexp"

	"github.com/51st-state/api/pkg/rbac"
)

var idRegexp = regexp.MustCompile(`^([a-z0-9-_]+\/)?[a-z0-9-_]+$`)

// Manager for managing role informations
//go:generate counterfeiter -o ./mocks/manager.go . Manager
type Manager interface {
	Get(context.Context, Identifier) (Complete, error)
	Set(context.Context, Complete) error
	Create(context.Context, Complete) error
	Delete(context.Context, Identifier) error
}

type manager struct {
	repository Repository
	rbac       rbac.Control
}

// NewManager creates a new instance of a manager for role information
func NewManager(repo Repository, rb rbac.Control) Manager {
	return &manager{
		repo,
		rb,
	}
}

var (
	errInvalidID = errors.New("invalid id format")
)

// Get role information
func (m *manager) Get(ctx context.Context, id Identifier) (Complete, error) {
	if !idRegexp.MatchString(string(id.ID())) {
		return nil, errInvalidID
	}

	c, err := m.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	rules, err := m.rbac.GetRoleRules(ctx, id.ID())
	if err != nil {
		return nil, err
	}

	c.Data().SetRules(rules)

	return c, nil
}

// Set role information
func (m *manager) Set(ctx context.Context, c Complete) error {
	if !idRegexp.MatchString(string(c.ID())) {
		return errInvalidID
	}

	if err := m.repository.Update(ctx, c); err != nil {
		return err
	}

	return m.rbac.SetRoleRules(ctx, c.ID(), c.Data().Rules)
}

// Create a role with role information
func (m *manager) Create(ctx context.Context, c Complete) error {
	if !idRegexp.MatchString(string(c.ID())) {
		return errInvalidID
	}

	if err := m.repository.Create(ctx, c); err != nil {
		return err
	}

	return m.rbac.SetRoleRules(ctx, c.ID(), c.Data().Rules)
}

// Delete role information
func (m *manager) Delete(ctx context.Context, id Identifier) error {
	if !idRegexp.MatchString(string(id.ID())) {
		return errInvalidID
	}

	if err := m.rbac.SetRoleRules(ctx, id.ID(), make(rbac.RoleRules, 0)); err != nil {
		return err
	}

	return m.repository.Delete(ctx, id)
}
