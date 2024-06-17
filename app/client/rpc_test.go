package client

import (
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// client

type ServiceClient interface {
	Get(addr string) RpcClient
}

type DefaultClient struct{}

func (c *DefaultClient) getOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func (c *DefaultClient) Get(addr string) RpcClient {
	return NewRpcClient(addr, c.getOptions()...)
}

// xx-service

type UserClient struct {
	DefaultClient
}

var (
	userAddr = "127.0.0.1:8080"
	client   RpcClient
	once     sync.Once
)

func GetUserClient() RpcClient {
	once.Do(func() {
		c := &UserClient{}
		client = c.Get(userAddr)
	})
	return client
}

func TestRpcClient(t *testing.T) {
	client := GetUserClient()
	conn := client.Get()
	defer client.Put(conn)

	// do something
}
