package rbac

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/51st-state/api/pkg/rbac/proto"
)

type grpcClient struct {
	client pb.ControlClient
}

// NewGRPCClient for the rbac controller
func NewGRPCClient(c *grpc.ClientConn) Control {
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

// GetAccountRoles returns the account roles
func (c *grpcClient) GetAccountRoles(ctx context.Context, accountID AccountID) (AccountRoles, error) {
	grpcRoles, err := c.client.GetAccountRoles(ctx, &pb.AccountID{
		ID: string(accountID),
	})
	if err != nil {
		return nil, err
	}

	roles := make(AccountRoles, 0)
	for _, v := range grpcRoles.GetRoleIDs() {
		roles = append(roles, RoleID(v))
	}

	return roles, nil
}

// SetAccountRoles sets the roles of a account
func (c *grpcClient) SetAccountRoles(ctx context.Context, accountID AccountID, roles AccountRoles) error {
	grpcRoles := &pb.AccountRoles{
		RoleIDs: []string{},
	}
	for _, v := range roles {
		grpcRoles.RoleIDs = append(grpcRoles.RoleIDs, string(v))
	}

	_, err := c.client.SetAccountRoles(ctx, &pb.SetAccountRolesRequest{
		AccountID: &pb.AccountID{
			ID: string(accountID),
		},
		AccountRoles: grpcRoles,
	})
	return err
}

// IsAccountAllowed checks whether a account has access to a rule
func (c *grpcClient) IsAccountAllowed(ctx context.Context, accountID AccountID, rule Rule) (bool, error) {
	resp, err := c.client.IsAccountAllowed(ctx, &pb.IsAccountAllowedRequest{
		AccountID: &pb.AccountID{
			ID: string(accountID),
		},
		Rule: &pb.Rule{
			Rule: string(rule),
		},
	})
	if err != nil {
		return false, err
	}

	return resp.GetAllowed(), nil
}
