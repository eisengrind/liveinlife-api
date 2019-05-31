package inventory

import (
	"context"

	pb "github.com/51st-state/api/pkg/apis/inventory/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type grpcServer struct {
	manager Manager
}

// NewGRPCServer creates a new grpc server instance
func NewGRPCServer(m Manager) pb.ManagerServer {
	return &grpcServer{
		m,
	}
}

func (s *grpcServer) Get(ctx context.Context, id *pb.Identifier) (*pb.Complete, error) {
	c, err := s.manager.Get(ctx, &identifier{id.GetGUID()})
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Item, 0)
	for _, v := range c.Data().Items {
		items = append(items, &pb.Item{
			ID:     v.ID,
			Amount: v.Amount,
			Subset: v.Subset,
		})
	}

	return &pb.Complete{
		Identifier: &pb.Identifier{
			GUID: c.GUID(),
		},
		Incomplete: &pb.Incomplete{
			Items: items,
		},
	}, nil
}

func (s *grpcServer) Create(ctx context.Context, inc *pb.Incomplete) (*pb.Complete, error) {
	items := make([]*Item, 0)
	for _, v := range inc.GetItems() {
		items = append(items, &Item{
			ID:     v.GetID(),
			Amount: v.GetAmount(),
			Subset: v.GetSubset(),
		})
	}

	c, err := s.manager.Create(ctx, NewIncomplete(items))
	if err != nil {
		return nil, err
	}

	return &pb.Complete{
		Identifier: &pb.Identifier{
			GUID: c.GUID(),
		},
		Incomplete: &pb.Incomplete{
			Items: inc.GetItems(),
		},
	}, nil
}

func (s *grpcServer) AddItem(ctx context.Context, req *pb.AddItemRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.manager.AddItem(
		ctx,
		&identifier{req.GetIdentifier().GetGUID()},
		&Item{
			ID:     req.GetItem().GetID(),
			Amount: req.GetItem().GetAmount(),
			Subset: req.GetItem().GetSubset(),
		},
	)
}

func (s *grpcServer) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.manager.RemoveItem(
		ctx,
		&identifier{req.GetIdentifier().GetGUID()},
		&Item{
			ID:     req.GetItem().GetID(),
			Amount: req.GetItem().GetAmount(),
			Subset: req.GetItem().GetSubset(),
		},
	)
}

func (s *grpcServer) Delete(ctx context.Context, id *pb.Identifier) (*empty.Empty, error) {
	return &empty.Empty{}, s.manager.Delete(
		ctx,
		&identifier{id.GetGUID()},
	)
}
