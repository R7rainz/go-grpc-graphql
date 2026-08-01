package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/r7rainz/go-grpc-graphql/catalog/internal/domain"
	"github.com/r7rainz/go-grpc-graphql/catalog/internal/service"
	"github.com/r7rainz/go-grpc-graphql/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service service.Service
}

func ListenGRPC(s service.Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(serv, &grpcServer{
		service: s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.PostProductResponse{Product: toProtoProduct(p)}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.GetProductResponse{Product: toProtoProduct(p)}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []domain.Product
	var err error

	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, int(r.Skip), int(r.Take))
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, int(r.Skip), int(r.Take))
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := make([]*pb.Product, 0, len(res))
	for i := range res {
		products = append(products, toProtoProduct(&res[i]))
	}

	return &pb.GetProductsResponse{Products: products}, nil
}

func toProtoProduct(p *domain.Product) *pb.Product {
	if p == nil {
		return nil
	}

	return &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}
}
