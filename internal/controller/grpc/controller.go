package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/ximura/gowallet/api"
	"github.com/ximura/gowallet/internal/core/domain"
	"github.com/ximura/gowallet/internal/core/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	service ports.WalletService

	api.UnimplementedWalletServiceServer
}

func NewWalletController(service ports.WalletService) api.WalletServiceServer {
	return server{
		service: service,
	}
}

func (s server) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{
		Message: "pong:" + req.Message,
	}, nil
}

func (s server) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	if req.Currency == "" {
		return nil, status.Errorf(codes.InvalidArgument, "currency can't be empty")
	}

	u, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "account id should be uuid")
	}
	w, err := s.service.Create(ctx, u, domain.Currency(req.Currency))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &api.CreateResponse{
		Wallet: convertWallet(w),
	}, nil
}

func (s server) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	u, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "account id should be uuid")
	}

	r, err := s.service.List(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var response api.ListResponse
	for i := range r {
		response.Wallet = append(response.Wallet, convertWallet(r[i]))
	}

	return &response, nil
}

func (s server) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	w, err := s.service.Get(ctx, int(req.WalletID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &api.GetResponse{
		Wallet: convertWallet(w),
	}, nil
}

func (s server) ProcessTransaction(ctx context.Context, req *api.Transaction) (*api.Wallet, error) {
	u, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id should be uuid")
	}

	w, err := s.service.ProcessTransaction(ctx, domain.Transaction{
		ID:       u,
		WalletID: int(req.WalletID),
		Amount:   int(req.Amount),
		Currency: domain.Currency(req.Currency),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return convertWallet(w), nil
}

func convertWallet(w domain.Wallet) *api.Wallet {
	return &api.Wallet{
		Id:       int32(w.ID),
		Customer: w.Account.String(),
		Amount:   int64(w.Amount),
		Currency: string(w.Currency),
	}
}
