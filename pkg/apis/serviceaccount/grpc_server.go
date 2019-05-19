package serviceaccount

import (
	"context"

	pb "github.com/51st-state/api/pkg/apis/serviceaccount/proto"
	"github.com/51st-state/api/pkg/rbac"
	proto1 "github.com/51st-state/api/pkg/rbac/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type grpcServer struct {
	manager Manager
}

// NewGRPCServer creates a new instance of a grpc server for managing service accounts
func NewGRPCServer(m Manager) pb.ManagerServer {
	return &grpcServer{m}
}

func (g *grpcServer) Get(ctx context.Context, id *pb.Identifier) (*pb.Complete, error) {
	c, err := g.manager.Get(ctx, &identifier{id.GetGUID()})
	if err != nil {
		return nil, err
	}

	return &pb.Complete{
		Identifier: id,
		Incomplete: &pb.Incomplete{
			Name:        c.Data().Name,
			Description: c.Data().Description,
		},
	}, nil
}

func (g *grpcServer) Create(ctx context.Context, inc *pb.Incomplete) (*pb.Complete, error) {
	c, err := g.manager.Create(ctx, NewIncomplete(inc.GetName(), inc.GetDescription()))
	if err != nil {
		return nil, err
	}

	return &pb.Complete{
		Identifier: &pb.Identifier{
			GUID: c.GUID(),
		},
		Incomplete: &pb.Incomplete{
			Name:        c.Data().Name,
			Description: c.Data().Description,
		},
	}, nil
}

func (g *grpcServer) Update(ctx context.Context, c *pb.Complete) (*empty.Empty, error) {
	return &empty.Empty{}, g.manager.Update(ctx, &complete{
		&identifier{
			c.GetIdentifier().GetGUID(),
		},
		NewIncomplete(
			c.GetIncomplete().GetName(),
			c.GetIncomplete().GetDescription(),
		),
	})
}

func (g *grpcServer) Delete(ctx context.Context, id *pb.Identifier) (*empty.Empty, error) {
	return &empty.Empty{}, g.manager.Delete(ctx, &identifier{id.GetGUID()})
}

func (g *grpcServer) GetRoles(ctx context.Context, id *pb.Identifier) (*proto1.AccountRoles, error) {
	roles, err := g.manager.GetRoles(ctx, &identifier{id.GetGUID()})
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

func (g *grpcServer) SetRoles(ctx context.Context, req *pb.SetRolesRequest) (*empty.Empty, error) {
	roles := make(rbac.AccountRoles, 0)
	for _, v := range req.GetRoles().GetRoleIDs() {
		roles = append(roles, rbac.RoleID(v))
	}

	return &empty.Empty{}, g.manager.SetRoles(
		ctx,
		&identifier{req.GetIdentifier().GetGUID()},
		roles,
	)
}
