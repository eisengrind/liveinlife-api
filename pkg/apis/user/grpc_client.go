package user

import (
	"context"

	pb "github.com/51st-state/api/pkg/apis/user/proto"
	"github.com/51st-state/api/pkg/rbac"
	proto1 "github.com/51st-state/api/pkg/rbac/proto"
	grpc "google.golang.org/grpc"
)

// GRPCClient for the user package
type GRPCClient struct {
	client pb.ManagerClient
}

// NewGRPCClient for users
func NewGRPCClient(c *grpc.ClientConn) *GRPCClient {
	return &GRPCClient{
		pb.NewManagerClient(c),
	}
}

// Get an user object
func (cli *GRPCClient) Get(ctx context.Context, id Identifier) (Complete, error) {
	resp, err := cli.client.GetUser(ctx, &pb.GetUserRequest{
		UUID: &pb.UUID{
			UUID: id.UUID(),
		},
	})
	if err != nil {
		return nil, err
	}

	return newComplete(
		id,
		NewIncomplete(
			WCFUserID(resp.GetData().GetWCFUserID()),
			resp.GetData().GetGameHash(),
			resp.GetData().GetBanned(),
		),
	), nil
}

// GetByWCFUserID returns a user filtered by its wcf user id
func (cli *GRPCClient) GetByWCFUserID(ctx context.Context, wcfUserID WCFUserID) (Complete, error) {
	resp, err := cli.client.GetUserByWCFUserID(ctx, &pb.GetUserByWCFUserIDRequest{
		WCFUserID: uint64(wcfUserID),
	})
	if err != nil {
		return nil, err
	}

	return newComplete(
		newIdentifier(resp.UUID.GetUUID()),
		NewIncomplete(
			WCFUserID(resp.GetData().GetWCFUserID()),
			resp.GetData().GetGameHash(),
			resp.GetData().GetBanned(),
		),
	), nil
}

// GetByGameSerialHash returns a user filtered by its game serial hash
func (cli *GRPCClient) GetByGameSerialHash(ctx context.Context, hash string) (Complete, error) {
	resp, err := cli.client.GetUserByGameSerialHash(ctx, &pb.GetUserByGameSerialHashRequest{
		Hash: hash,
	})
	if err != nil {
		return nil, err
	}

	return newComplete(
		newIdentifier(resp.UUID.GetUUID()),
		NewIncomplete(
			WCFUserID(resp.GetData().GetWCFUserID()),
			resp.GetData().GetGameHash(),
			resp.GetData().GetBanned(),
		),
	), nil
}

// Create an user object
func (cli *GRPCClient) Create(ctx context.Context, inc Incomplete) (Complete, error) {
	resp, err := cli.client.CreateUser(ctx, &pb.CreateUserRequest{
		Data: &pb.Data{
			WCFUserID: uint64(inc.Data().WCFUserID),
			GameHash:  inc.Data().GameSerialHash,
			Banned:    inc.Data().Banned,
		},
	})
	if err != nil {
		return nil, err
	}

	return newComplete(
		newIdentifier(resp.GetUUID().GetUUID()),
		inc,
	), nil
}

// Update an user object
func (cli *GRPCClient) Update(ctx context.Context, c Complete) error {
	_, err := cli.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		UUID: &pb.UUID{
			UUID: c.UUID(),
		},
		Data: &pb.Data{
			WCFUserID: uint64(c.Data().WCFUserID),
			GameHash:  c.Data().GameSerialHash,
			Banned:    c.Data().Banned,
		},
	})

	return err
}

// Delete a user
func (cli *GRPCClient) Delete(ctx context.Context, id Identifier) error {
	_, err := cli.client.DeleteUser(ctx, &pb.DeleteUserRequest{
		UUID: &pb.UUID{
			UUID: id.UUID(),
		},
	})

	return err
}

// CheckPassword of an user
func (cli *GRPCClient) CheckPassword(ctx context.Context, id Identifier, pw IncompletePassword) error {
	_, err := cli.client.CheckUserPassword(ctx, &pb.CheckUserPasswordRequest{
		UUID: &pb.UUID{
			UUID: id.UUID(),
		},
		Password: &pb.IncompletePassword{
			Password: pw.Password(),
		},
	})

	return err
}

// GetWCFInfoByEmail of a wcf user
func (cli *GRPCClient) GetWCFInfoByEmail(ctx context.Context, email string) (*WCFUserInfo, error) {
	resp, err := cli.client.GetWCFInfoByEmail(ctx, &pb.GetWCFInfoByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return &WCFUserInfo{
		WCFUserID(resp.GetUserID()),
		resp.GetEmail(),
		newCompletePassword(resp.GetPassword()),
	}, nil
}

// GetWCFInfoByUsername of a wcf user
func (cli *GRPCClient) GetWCFInfoByUsername(ctx context.Context, username string) (*WCFUserInfo, error) {
	resp, err := cli.client.GetWCFInfoByUsername(ctx, &pb.GetWCFInfoByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	return &WCFUserInfo{
		WCFUserID(resp.GetUserID()),
		resp.GetEmail(),
		newCompletePassword(resp.GetPassword()),
	}, nil
}

// GetRoles of a user
func (cli *GRPCClient) GetRoles(ctx context.Context, id Identifier) (rbac.SubjectRoles, error) {
	resp, err := cli.client.GetUserRoles(ctx, &pb.UUID{
		UUID: id.UUID(),
	})
	if err != nil {
		return nil, err
	}

	roles := make(rbac.SubjectRoles, 0)
	for _, v := range resp.GetRoleIDs() {
		roles = append(roles, rbac.RoleID(v))
	}

	return nil, nil
}

// SetRoles of a user
func (cli *GRPCClient) SetRoles(ctx context.Context, id Identifier, roles rbac.SubjectRoles) error {
	grpcRoles := make([]string, 0)
	for _, v := range roles {
		grpcRoles = append(grpcRoles, string(v))
	}

	_, err := cli.client.SetUserRoles(ctx, &pb.SetUserRolesRequest{
		UUID: &pb.UUID{
			UUID: id.UUID(),
		},
		Roles: &proto1.SubjectRoles{
			RoleIDs: grpcRoles,
		},
	})
	return err
}
