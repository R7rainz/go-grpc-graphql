package client

import (
	"context"

	"github.com/r7rainz/go-grpc-graphql/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

type Account struct {
	ID   string
	Name string
}

func NewClient(target string) (*Client, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		service: pb.NewAccountServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	res, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return toAccount(res.GetAccount()), nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	res, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return toAccount(res.GetAccount()), nil
}

func (c *Client) GetAccounts(ctx context.Context, skip int, take int) ([]*Account, error) {
	res, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{
		Skip: int32(skip),
		Take: int32(take),
	})
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0, len(res.GetAccounts()))
	for _, account := range res.GetAccounts() {
		accounts = append(accounts, toAccount(account))
	}

	return accounts, nil
}

func toAccount(account *pb.Account) *Account {
	if account == nil {
		return nil
	}

	return &Account{
		ID:   account.GetId(),
		Name: account.GetName(),
	}
}
