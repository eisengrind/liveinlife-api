package rbac

import (
	"context"
	"errors"
)

// Control of the rbac system
type Control struct {
	repository Repository
}

// NewControl instantiates a new RBAC control
//go:generate protoc -I ./proto --go_out=plugins=grpc:./proto ./proto/control.proto
func NewControl(r Repository) *Control {
	return &Control{
		r,
	}
}

var (
	errEmptyRoleID    = errors.New("empty role id")
	errEmptyRule      = errors.New("empty rule")
	errEmptySubjectID = errors.New("empty subject id")
)

// GetRoleRules gets the rules of  role
func (m *Control) GetRoleRules(ctx context.Context, roleID RoleID) (RoleRules, error) {
	if roleID == "" {
		return nil, errEmptyRoleID
	}

	return m.GetRoleRules(ctx, roleID)
}

// SetRoleRules sets the rules of a role
func (m *Control) SetRoleRules(ctx context.Context, roleID RoleID, rules RoleRules) error {
	if roleID == "" {
		return errEmptyRoleID
	}

	for _, v := range rules {
		if v == "" {
			return errEmptyRule
		}
	}

	return nil
}

// GetSubjectRoles returns the subject roles
func (m *Control) GetSubjectRoles(ctx context.Context, subjectID SubjectID) (SubjectRoles, error) {
	if subjectID == "" {
		return nil, errEmptySubjectID
	}

	return m.repository.GetSubjectRoles(ctx, subjectID)
}

// SetSubjectRoles sets the roles of a subject
func (m *Control) SetSubjectRoles(ctx context.Context, subjectID SubjectID, roles SubjectRoles) error {
	if subjectID == "" {
		return errEmptySubjectID
	}

	for _, v := range roles {
		if v == "" {
			return errEmptyRoleID
		}
	}

	return nil
}

// IsSubjectAllowed checks whether a subject has access to a rule
func (m *Control) IsSubjectAllowed(ctx context.Context, subjectID SubjectID, rule Rule) error {
	if subjectID == "" {
		return errEmptySubjectID
	}

	if rule == "" {
		return errEmptyRule
	}

	count, err := m.repository.GetSubjectRuleCount(ctx, subjectID, rule)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("insufficient permissions")
	}

	return nil
}
