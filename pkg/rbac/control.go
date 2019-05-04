package rbac

import (
	"context"
	"errors"
)

// Control of the rbac system
//go:generate counterfeiter -o ./mocks/control.go . Control
type Control interface {
	GetRoleRules(ctx context.Context, roleID RoleID) (RoleRules, error)
	SetRoleRules(ctx context.Context, roleID RoleID, rules RoleRules) error
	GetAccountRoles(ctx context.Context, accountID AccountID) (AccountRoles, error)
	SetAccountRoles(ctx context.Context, accountID AccountID, roles AccountRoles) error
	IsAccountAllowed(ctx context.Context, accountID AccountID, rule Rule) error
}

type control struct {
	repository Repository
}

// NewControl instantiates a new RBAC control
//go:generate protoc -I ./proto --go_out=plugins=grpc:./proto ./proto/control.proto
func NewControl(r Repository) Control {
	return &control{
		r,
	}
}

var (
	errEmptyRoleID    = errors.New("empty role id")
	errEmptyRule      = errors.New("empty rule")
	errEmptyAccountID = errors.New("empty account id")
)

// GetRoleRules gets the rules of  role
func (m *control) GetRoleRules(ctx context.Context, roleID RoleID) (RoleRules, error) {
	if roleID == "" {
		return nil, errEmptyRoleID
	}

	return m.repository.GetRoleRules(ctx, roleID)
}

// SetRoleRules sets the rules of a role
func (m *control) SetRoleRules(ctx context.Context, roleID RoleID, rules RoleRules) error {
	if roleID == "" {
		return errEmptyRoleID
	}

	for _, v := range rules {
		if v == "" {
			return errEmptyRule
		}
	}

	return m.repository.SetRoleRules(ctx, roleID, rules)
}

// GetAccountRoles returns the account roles
func (m *control) GetAccountRoles(ctx context.Context, accountID AccountID) (AccountRoles, error) {
	if accountID == "" {
		return nil, errEmptyAccountID
	}

	return m.repository.GetAccountRoles(ctx, accountID)
}

// SetAccountRoles sets the roles of a account
func (m *control) SetAccountRoles(ctx context.Context, accountID AccountID, roles AccountRoles) error {
	if accountID == "" {
		return errEmptyAccountID
	}

	for _, v := range roles {
		if v == "" {
			return errEmptyRoleID
		}
	}

	return m.repository.SetAccountRoles(ctx, accountID, roles)
}

// IsAccountAllowed checks whether a account has access to a rule
func (m *control) IsAccountAllowed(ctx context.Context, accountID AccountID, rule Rule) error {
	if accountID == "" {
		return errEmptyAccountID
	}

	if rule == "" {
		return errEmptyRule
	}

	count, err := m.repository.GetAccountRuleCount(ctx, accountID, rule)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("insufficient permissions")
	}

	return nil
}
