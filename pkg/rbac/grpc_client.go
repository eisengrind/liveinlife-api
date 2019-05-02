package rbac

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/51st-state/api/pkg/rbac/proto"
)

type grpcClient struct {
	client pb.ControlClient
}

// NewgrpcClient for the rbac controller
func NewgrpcClient(c *grpc.ClientConn) Control {
	return &grpcClient{
		pb.NewControlClient(c),
	}
}

// GetRoleRules gets the rules of  role
func (c *grpcClient) GetRoleRules(ctx context.Context, roleID RoleID) (RoleRules, error) {
	grpcRules, err := c.client.GetRoleRules(ctx, &pb.RoleID{
		ID: string(roleID),
	})
	if err != nil {
		return nil, err
	}

	rules := make(RoleRules, 0)
	for _, v := range grpcRules.GetRules() {
		rules = append(rules, Rule(v))
	}

	return rules, nil
}

// SetRoleRules sets the rules of a role
func (c *grpcClient) SetRoleRules(ctx context.Context, roleID RoleID, rules RoleRules) error {
	grpcRules := &pb.RoleRules{
		Rules: []string{},
	}
	for _, v := range rules {
		grpcRules.Rules = append(grpcRules.Rules, string(v))
	}

	_, err := c.client.SetRoleRules(ctx, &pb.SetRoleRulesRequest{
		RoleID: &pb.RoleID{
			ID: string(roleID),
		},
		RoleRules: grpcRules,
	})
	return err
}

// GetSubjectRoles returns the subject roles
func (c *grpcClient) GetSubjectRoles(ctx context.Context, subjectID SubjectID) (SubjectRoles, error) {
	grpcRoles, err := c.client.GetSubjectRoles(ctx, &pb.SubjectID{
		ID: string(subjectID),
	})
	if err != nil {
		return nil, err
	}

	roles := make(SubjectRoles, 0)
	for _, v := range grpcRoles.GetRoleIDs() {
		roles = append(roles, RoleID(v))
	}

	return roles, nil
}

// SetSubjectRoles sets the roles of a subject
func (c *grpcClient) SetSubjectRoles(ctx context.Context, subjectID SubjectID, roles SubjectRoles) error {
	grpcRoles := &pb.SubjectRoles{
		RoleIDs: []string{},
	}
	for _, v := range roles {
		grpcRoles.RoleIDs = append(grpcRoles.RoleIDs, string(v))
	}

	_, err := c.client.SetSubjectRoles(ctx, &pb.SetSubjectRolesRequest{
		SubjectID: &pb.SubjectID{
			ID: string(subjectID),
		},
		SubjectRoles: grpcRoles,
	})
	return err
}

// IsSubjectAllowed checks whether a subject has access to a rule
func (c *grpcClient) IsSubjectAllowed(ctx context.Context, subjectID SubjectID, rule Rule) error {
	_, err := c.client.IsSubjectAllowed(ctx, &pb.IsSubjectAllowedRequest{
		SubjectID: &pb.SubjectID{
			ID: string(subjectID),
		},
		Rule: &pb.Rule{
			Rule: string(rule),
		},
	})
	return err
}
