package user

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/51st-state/api/pkg/apis/user/proto"
	"github.com/51st-state/api/pkg/rbac"
	proto1 "github.com/51st-state/api/pkg/rbac/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

// GRPCServer for external user management
type GRPCServer struct {
	manager Manager
}

// NewGRPCServer for user management
func NewGRPCServer(m Manager) pb.ManagerServer {
	return &GRPCServer{m}
}

// GetUser from the user database
func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	c, err := s.manager.Get(ctx, newIdentifier(req.GetUUID().GetUUID()))
	if err != nil {
		return nil, err
	}

	return &pb.User{
		UUID: &pb.UUID{
			UUID: c.UUID(),
		},
		Data: &pb.Data{
			WCFUserID: uint64(c.Data().WCFUserID),
			GameHash:  c.Data().GameSerialHash,
			Banned:    c.Data().Banned,
		},
	}, nil
}

// GetUserByWCFUserID returns a user filtered by its wcf user id
func (s *GRPCServer) GetUserByWCFUserID(ctx context.Context, req *pb.GetUserByWCFUserIDRequest) (*pb.User, error) {
	c, err := s.manager.GetByWCFUserID(ctx, WCFUserID(req.GetWCFUserID()))
	if err != nil {
		return nil, err
	}

	return &pb.User{
		UUID: &pb.UUID{
			UUID: c.UUID(),
		},
		Data: &pb.Data{
			WCFUserID: uint64(c.Data().WCFUserID),
			GameHash:  c.Data().GameSerialHash,
			Banned:    c.Data().Banned,
		},
	}, nil
}

// GetUserByGameSerialHash returns a user filtered by its game serial hash
func (s *GRPCServer) GetUserByGameSerialHash(ctx context.Context, req *pb.GetUserByGameSerialHashRequest) (*pb.User, error) {
	c, err := s.manager.GetByGameSerialHash(ctx, req.GetHash())
	if err == sql.ErrNoRows {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	} else if err != nil {
		return nil, err
	}

	return &pb.User{
		UUID: &pb.UUID{
			UUID: c.UUID(),
		},
		Data: &pb.Data{
			WCFUserID: uint64(c.Data().WCFUserID),
			GameHash:  c.Data().GameSerialHash,
			Banned:    c.Data().Banned,
		},
	}, nil
}

// CreateUser in the database
func (s *GRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	c, err := s.manager.Create(ctx, NewIncomplete(
		WCFUserID(req.GetData().GetWCFUserID()),
		req.GetData().GetGameHash(),
		req.GetData().GetBanned(),
	))
	if err != nil {
		return nil, err
	}

	return &pb.User{
		UUID: &pb.UUID{
			UUID: c.UUID(),
		},
		Data: &pb.Data{
			WCFUserID: uint64(c.Data().WCFUserID),
			GameHash:  c.Data().GameSerialHash,
			Banned:    c.Data().Banned,
		},
	}, nil
}

// DeleteUser from the database
func (s *GRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.manager.Delete(
		ctx,
		newIdentifier(
			req.GetUUID().GetUUID(),
		),
	)
}

// UpdateUser credentials in the database
func (s *GRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.manager.Update(ctx, newComplete(
		newIdentifier(
			req.GetUUID().GetUUID(),
		),
		NewIncomplete(
			WCFUserID(req.GetData().GetWCFUserID()),
			req.GetData().GetGameHash(),
			req.GetData().GetBanned(),
		),
	))
}

// CheckUserPassword if the given password matches
func (s *GRPCServer) CheckUserPassword(ctx context.Context, req *pb.CheckUserPasswordRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.manager.CheckPassword(
		ctx,
		newIdentifier(
			req.GetUUID().GetUUID(),
		),
		newIncompletePassword(
			req.GetPassword().GetPassword(),
		),
	)
}

// GetWCFInfo of a wcf user
func (s *GRPCServer) GetWCFInfo(ctx context.Context, req *pb.GetWCFInfoRequest) (*pb.WCFUserInfo, error) {
	info, err := s.manager.GetWCFInfo(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	return &pb.WCFUserInfo{
		UserID:   uint64(info.UserID),
		Email:    info.Email,
		Password: info.Password.Hash(),
	}, nil
}

// GetUserRoles of a user
func (s *GRPCServer) GetUserRoles(ctx context.Context, id *pb.UUID) (*proto1.AccountRoles, error) {
	roles, err := s.manager.GetRoles(ctx, newIdentifier(id.GetUUID()))
	if err != nil {
		return nil, err
	}

	grpcRoles := &proto1.AccountRoles{
		RoleIDs: make([]string, 0),
	}
	for _, v := range roles {
		grpcRoles.RoleIDs = append(grpcRoles.RoleIDs, string(v))
	}

	return grpcRoles, nil
}

// SetUserRoles of a user
func (s *GRPCServer) SetUserRoles(ctx context.Context, req *pb.SetUserRolesRequest) (*empty.Empty, error) {
	roles := make(rbac.AccountRoles, 0)
	for _, v := range req.GetRoles().GetRoleIDs() {
		roles = append(roles, rbac.RoleID(v))
	}

	return &empty.Empty{}, s.manager.SetRoles(ctx, newIdentifier(req.GetUUID().GetUUID()), roles)
}
