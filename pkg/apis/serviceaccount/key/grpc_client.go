package key

import (
	"context"
	"crypto/rsa"
	"encoding/json"

	pb "github.com/51st-state/api/pkg/apis/serviceaccount/key/proto"
	"google.golang.org/grpc"
)

type grpcClient struct {
	client pb.ManagerClient
}

// NewGRPCClient creates a new grpc client for serviceaccount keys
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

	var publicKey rsa.PublicKey
	if err := json.Unmarshal(resp.GetIncomplete().GetPublicKey(), &publicKey); err != nil {
		return nil, err
	}

	inc := NewIncomplete(
		resp.GetIncomplete().Name,
		resp.GetIncomplete().Description,
	)
	inc.Data().PublicKey = &publicKey

	return &complete{
		&identifier{
			resp.GetIdentifier().GetGUID(),
		},
		inc,
	}, nil
}

func (g *grpcClient) Update(ctx context.Context, c Complete) error {
	b, err := json.Marshal(c.Data().PublicKey)
	if err != nil {
		return err
	}

	_, err = g.client.Update(ctx, &pb.Complete{
		Identifier: &pb.Identifier{
			GUID: c.GUID(),
		},
		Incomplete: &pb.Incomplete{
			Name:               c.Data().Name,
			Description:        c.Data().Description,
			PublicKey:          b,
			ServiceAccountGUID: c.Data().ServiceAccountGUID,
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

func (g *grpcClient) Create(ctx context.Context, inc Incomplete) (*ClientKey, error) {
	resp, err := g.client.Create(ctx, &pb.Incomplete{
		Name:               inc.Data().Name,
		Description:        inc.Data().Description,
		ServiceAccountGUID: inc.Data().ServiceAccountGUID,
	})
	if err != nil {
		return nil, err
	}

	var privateKey rsa.PrivateKey
	if err := json.Unmarshal(resp.GetPrivateKey(), &privateKey); err != nil {
		return nil, err
	}

	return &ClientKey{
		ServiceAccountGUID: resp.GetServiceAccountGUID(),
		GUID:               resp.GetGUID(),
		PrivateKey:         &privateKey,
	}, nil
}
