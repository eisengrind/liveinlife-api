package serviceaccount

import (
	"context"

	pb "github.com/51st-state/api/pkg/apis/serviceaccount/proto"
	"github.com/51st-state/api/pkg/rbac"
	proto1 "github.com/51st-state/api/pkg/rbac/proto"
	"google.golang.org/grpc"
)

type grpcClient struct {
	client pb.ManagerClient
}

// NewGRPCClient creates a new grpc client for the service account manager
func NewGRPCClient(c *grpc.ClientConn) Manager {
	return &grpcClient{
		pb.NewManagerClient(c),
	}
}

func (g *grpcClient) Get(ctx context.Context, id Identifier) (Complete, error) {
	resp, err := g.client.Get(ctx, &pb.Identifier{
		GUID: id.GUID(),
	})
	if err != nil {
		return nil, err
	}

	return &complete{
		&identifier{
			resp.GetIdentifier().GetGUID(),
		},
		NewIncomplete(
			resp.GetIncomplete().GetName(),
			resp.GetIncomplete().GetDescription(),
		),
	}, nil
}

func (g *grpcClient) Create(ctx context.Context, inc Incomplete) (Complete, error) {
	resp, err := g.client.Create(ctx, &pb.Incomplete{
		Name:        inc.Data().Name,
		Description: inc.Data().Description,
	})
	if err != nil {
		return nil, err
	}

	return &complete{
		&identifier{
			resp.GetIdentifier().GetGUID(),
		},
		inc,
	}, nil
}

func (g *grpcClient) Update(ctx context.Context, c Complete) error {
	_, err := g.client.Update(ctx, &pb.Complete{
		Identifier: &pb.Identifier{
			GUID: c.GUID(),
		},
		Incomplete: &pb.Incomplete{
			Name:        c.Data().Name,
			Description: c.Data().Description,
		},
	})
	return err
}

func (g *grpcClient) Delete(ctx context.Context, id Identifier) error {
	_, err := g.client.Delete(ctx, &pb.Identifier{
		GUID: id.GUID(),
	})
	return err
}

func (g *grpcClient) GetRoles(ctx context.Context, id Identifier) (rbac.AccountRoles, error) {
	resp, err := g.client.GetRoles(ctx, &pb.Identifier{
		GUID: id.GUID(),
	})
	if err != nil {
		return nil, err
	}

	roles := make(rbac.AccountRoles, 0)
	for _, v := range resp.GetRoleIDs() {
		roles = append(roles, rbac.RoleID(v))
	}

	return roles, nil
}

func (g *grpcClient) SetRoles(ctx context.Context, id Identifier, roles rbac.AccountRoles) error {
	grpcRoles := make([]string, 0)
	for _, v := range roles {
		grpcRoles = append(grpcRoles, string(v))
	}

	_, err := g.client.SetRoles(ctx, &pb.SetRolesRequest{
		Identifier: &pb.Identifier{
			GUID: id.GUID(),
		},
		Roles: &proto1.AccountRoles{
			RoleIDs: grpcRoles,
		},
	})
	return err
}
