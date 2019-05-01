package rbac

import "context"

// Repository for persistent RBAC storage
//go:generate counterfeiter -o ./mocks/repository.go . Repository
type Repository interface {
	// GetRoleRules fetches all available Rules from a role
	GetRoleRules(context.Context, RoleID) (RoleRules, error)
	// SetRoleRules sets the rules of a role
	SetRoleRules(context.Context, RoleID, RoleRules) error
	// GetSubjectRoles returns the roles of a subject
	GetSubjectRoles(context.Context, SubjectID) (SubjectRoles, error)
	// SetSubjectRoles sets the roles of a subject
	SetSubjectRoles(context.Context, SubjectID, SubjectRoles) error
	// GetSubjectRuleCount returns the amount of occurrences of a given rule
	// for a given subject
	GetSubjectRuleCount(context.Context, SubjectID, Rule) (uint64, error)
}
