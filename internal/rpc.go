package wallet

import (
	"context"

	"github.com/ximura/gowallet/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCService struct {
}

func NewGRPCService() *GRPCService {
	return &GRPCService{}
}

// StartServer starts a RPC server
func (gs *GRPCService) StartServer(ctx context.Context) (*grpc.Server, error) {

	s := grpc.NewServer()
	api.RegisterWalletServiceServer(s, api.UnimplementedWalletServiceServer{})
	reflection.Register(s)

	return s, nil
}
