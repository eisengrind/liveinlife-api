package rbac

import (
	"context"

	pb "github.com/51st-state/api/pkg/rbac/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type grpcServer struct {
	control Control
}

// NewGRPCServer creates a new grpc server instance for the rbac control
func NewGRPCServer(c Control) pb.ControlServer {
	return &grpcServer{
		c,
	}
}

func (s *grpcServer) GetRoleRules(ctx context.Context, roleID *pb.RoleID) (*pb.RoleRules, error) {
	roleRules, err := s.control.GetRoleRules(ctx, RoleID(roleID.GetID()))
	if err != nil {
		return nil, err
	}

	grpcRoleRules := make([]string, 0)
	for _, v := range roleRules {
		grpcRoleRules = append(grpcRoleRules, string(v))
	}

	return &pb.RoleRules{
		Rules: grpcRoleRules,
	}, nil
}

func (s *grpcServer) SetRoleRules(ctx context.Context, req *pb.SetRoleRulesRequest) (*empty.Empty, error) {
	roleRules := make(RoleRules, 0)
	for _, v := range req.GetRoleRules().GetRules() {
		roleRules = append(roleRules, Rule(v))
	}

	return &empty.Empty{}, s.control.SetRoleRules(ctx, RoleID(req.GetRoleID().GetID()), roleRules)
}

func (s *grpcServer) GetAccountRoles(ctx context.Context, accountID *pb.AccountID) (*pb.AccountRoles, error) {
	accountRoles, err := s.control.GetAccountRoles(ctx, AccountID(accountID.GetID()))
	if err != nil {
		return nil, err
	}

	grpcAccountRoles := make([]string, 0)
	for _, v := range accountRoles {
		grpcAccountRoles = append(grpcAccountRoles, string(v))
	}

	return &pb.AccountRoles{
		RoleIDs: grpcAccountRoles,
	}, nil
}

func (s *grpcServer) SetAccountRoles(ctx context.Context, req *pb.SetAccountRolesRequest) (*empty.Empty, error) {
	accountRoles := make(AccountRoles, 0)
	for _, v := range req.GetAccountRoles().GetRoleIDs() {
		accountRoles = append(accountRoles, RoleID(v))
	}

	return &empty.Empty{}, s.control.SetAccountRoles(ctx, AccountID(req.GetAccountID().GetID()), accountRoles)
}

func (s *grpcServer) IsAccountAllowed(ctx context.Context, req *pb.IsAccountAllowedRequest) (*pb.IsAccountAllowedResponse, error) {
	allowed, err := s.control.IsAccountAllowed(ctx, AccountID(req.GetAccountID().GetID()), Rule(req.GetRule().GetRule()))
	if err != nil {
		return nil, err
	}

	return &pb.IsAccountAllowedResponse{
		Allowed: allowed,
	}, nil
}
