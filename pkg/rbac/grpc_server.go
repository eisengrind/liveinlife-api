package rbac

import (
	"context"

	pb "github.com/51st-state/api/pkg/rbac/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type grpcServer struct {
	control *Control
}

// NewGRPCServer creates a new grpc server instance for the rbac control
func NewGRPCServer(c *Control) pb.ControlServer {
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

func (s *grpcServer) GetSubjectRoles(ctx context.Context, subjectID *pb.SubjectID) (*pb.SubjectRoles, error) {
	subjectRoles, err := s.control.GetSubjectRoles(ctx, SubjectID(subjectID.GetID()))
	if err != nil {
		return nil, err
	}

	grpcSubjectRoles := make([]string, 0)
	for _, v := range subjectRoles {
		grpcSubjectRoles = append(grpcSubjectRoles, string(v))
	}

	return &pb.SubjectRoles{
		RoleIDs: grpcSubjectRoles,
	}, nil
}

func (s *grpcServer) SetSubjectRoles(ctx context.Context, req *pb.SetSubjectRolesRequest) (*empty.Empty, error) {
	subjectRoles := make(SubjectRoles, 0)
	for _, v := range req.GetSubjectRoles().GetRoleIDs() {
		subjectRoles = append(subjectRoles, RoleID(v))
	}

	return &empty.Empty{}, s.control.SetSubjectRoles(ctx, SubjectID(req.GetSubjectID().GetID()), subjectRoles)
}

func (s *grpcServer) IsSubjectAllowed(ctx context.Context, req *pb.IsSubjectAllowedRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.control.IsSubjectAllowed(ctx, SubjectID(req.GetSubjectID().GetID()), Rule(req.GetRule().GetRule()))
}
