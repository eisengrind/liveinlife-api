package inventory

import (
	"context"

	pb "github.com/51st-state/api/pkg/apis/inventory/proto"
	"google.golang.org/grpc"
)

type grpcClient struct {
	client pb.ManagerClient
}

// NewGRPCClient creates a new inventory manager with a grpc conn
func NewGRPCClient(c *grpc.ClientConn) Manager {
	return &grpcClient{
		pb.NewManagerClient(c),
	}
}

func (g *grpcClient) Get(ctx context.Context, id Identifier) (Complete, error) {
	c, err := g.client.Get(ctx, &pb.Identifier{
		GUID: id.GUID(),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*Item, 0)
	for _, v := range c.GetIncomplete().GetItems() {
		items = append(items, &Item{
			ID:     v.GetID(),
			Amount: v.GetAmount(),
			Subset: v.GetSubset(),
		})
	}

	return &complete{
		id,
		NewIncomplete(items),
	}, nil
}

func (g *grpcClient) Create(ctx context.Context, inc Incomplete) (Complete, error) {
	items := make([]*pb.Item, 0)
	for _, v := range inc.Data().Items {
		items = append(items, &pb.Item{
			ID:     v.ID,
			Amount: v.Amount,
			Subset: v.Subset,
		})
	}
	c, err := g.client.Create(ctx, &pb.Incomplete{
		Items: items,
	})
	if err != nil {
		return nil, err
	}

	return &complete{
		&identifier{c.GetIdentifier().GetGUID()},
		inc,
	}, nil
}

func (g *grpcClient) AddItem(ctx context.Context, id Identifier, item *Item) error {
	_, err := g.client.AddItem(ctx, &pb.AddItemRequest{
		Identifier: &pb.Identifier{
			GUID: id.GUID(),
		},
		Item: &pb.Item{
			ID:     item.ID,
			Amount: item.Amount,
			Subset: item.Subset,
		},
	})
	return err
}

func (g *grpcClient) RemoveItem(ctx context.Context, id Identifier, item *Item) error {
	_, err := g.client.RemoveItem(ctx, &pb.RemoveItemRequest{
		Identifier: &pb.Identifier{
			GUID: id.GUID(),
		},
		Item: &pb.Item{
			ID:     item.ID,
			Amount: item.Amount,
			Subset: item.Subset,
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
