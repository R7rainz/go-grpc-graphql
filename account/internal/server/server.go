package server

import (
	"context"
	"fmt"
	"net"

	"github.com/r7rainz/go-grpc-graphql/account/internal/domain"
	"github.com/r7rainz/go-grpc-graphql/account/internal/service"
	"github.com/r7rainz/go-grpc-graphql/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service service.Service
}

func ListenGRPC(s service.Service, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcSrv, &grpcServer{
		service: s,
	})
	reflection.Register(grpcSrv)

	return grpcSrv.Serve(listener)
}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	account, err := s.service.PostAccount(ctx, r.GetName())
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{
		Account: toProtoAccount(account),
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.service.GetAccount(ctx, r.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Account: toProtoAccount(account),
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := s.service.GetAccounts(ctx, int(r.GetSkip()), int(r.GetTake()))
	if err != nil {
		return nil, err
	}

	protoAccounts := make([]*pb.Account, 0, len(accounts))
	for i := range accounts {
		protoAccounts = append(protoAccounts, toProtoAccount(&accounts[i]))
	}

	return &pb.GetAccountsResponse{
		Accounts: protoAccounts,
	}, nil
}

func toProtoAccount(account *domain.Account) *pb.Account {
	if account == nil {
		return nil
	}

	return &pb.Account{
		Id:   account.ID,
		Name: account.Name,
	}
}
