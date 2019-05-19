package key

import (
	"context"
	"crypto/rsa"
	"encoding/json"

	pb "github.com/51st-state/api/pkg/apis/serviceaccount/key/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type grpcServer struct {
	manager Manager
}

// NewGRPCServer creates a new instance of a grpc server interface to register
func NewGRPCServer(m Manager) pb.ManagerServer {
	return &grpcServer{m}
}

func (s *grpcServer) Get(ctx context.Context, id *pb.Identifier) (*pb.Complete, error) {
	c, err := s.manager.Get(ctx, &identifier{
		id.GetGUID(),
	})
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(c.Data().PublicKey)
	if err != nil {
		return nil, err
	}

	return &pb.Complete{
		Identifier: id,
		Incomplete: &pb.Incomplete{
			Name:               c.Data().Name,
			Description:        c.Data().Description,
			PublicKey:          b,
			ServiceAccountGUID: c.Data().ServiceAccountGUID,
		},
	}, nil
}

func (s *grpcServer) Create(ctx context.Context, req *pb.Incomplete) (*pb.ClientKey, error) {
	inc := NewIncomplete(
		req.GetName(),
		req.GetDescription(),
	)
	inc.Data().ServiceAccountGUID = req.GetServiceAccountGUID()

	clientKey, err := s.manager.Create(
		ctx,
		inc,
	)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(clientKey.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &pb.ClientKey{
		GUID:               clientKey.GUID,
		ServiceAccountGUID: clientKey.ServiceAccountGUID,
		PrivateKey:         b,
	}, nil
}

func (s *grpcServer) Update(ctx context.Context, c *pb.Complete) (*empty.Empty, error) {
	var publicKey rsa.PublicKey
	if err := json.Unmarshal(c.GetIncomplete().GetPublicKey(), &publicKey); err != nil {
		return nil, err
	}

	inc := NewIncomplete(
		c.GetIncomplete().GetName(),
		c.GetIncomplete().GetDescription(),
	)
	inc.Data().PublicKey = &publicKey

	return &empty.Empty{}, s.manager.Update(ctx, &complete{
		&identifier{
			c.GetIdentifier().GetGUID(),
		},
		inc,
	})
}

func (s *grpcServer) Delete(ctx context.Context, id *pb.Identifier) (*empty.Empty, error) {
	return &empty.Empty{}, s.manager.Delete(ctx, &identifier{
		id.GetGUID(),
	})
}
