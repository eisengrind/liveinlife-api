package rbac

import "context"

// Repository for persistent RBAC storage
//go:generate counterfeiter -o ./mocks/repository.go . Repository
type Repository interface {
	// GetRoleRules fetches all available Rules from a role
	GetRoleRules(context.Context, RoleID) (RoleRules, error)
	// SetRoleRules sets the rules of a role
	SetRoleRules(context.Context, RoleID, RoleRules) error
	// GetAccountRoles returns the roles of a subject
	GetAccountRoles(context.Context, AccountID) (AccountRoles, error)
	// SetAccountRoles sets the roles of a subject
	SetAccountRoles(context.Context, AccountID, AccountRoles) error
	// GetAccountRuleCount returns the amount of occurrences of a given rule
	// for a given subject
	GetAccountRuleCount(context.Context, AccountID, Rule) (uint64, error)
}
