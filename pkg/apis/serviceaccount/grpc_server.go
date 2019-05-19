package serviceaccount

import (
	"context"

	pb "github.com/51st-state/api/pkg/apis/serviceaccount/proto"
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
