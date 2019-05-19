package serviceaccount

import (
	"context"

	pb "github.com/51st-state/api/pkg/apis/serviceaccount/proto"
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
