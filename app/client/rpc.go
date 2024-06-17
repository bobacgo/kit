package client

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type RpcClient interface {
	Get() *grpc.ClientConn
	Put(*grpc.ClientConn)
}

type rpcClient struct {
	pool *sync.Pool
}

func NewRpcClient(target string, opts ...grpc.DialOption) RpcClient {
	return &rpcClient{pool: &sync.Pool{New: func() any {
		conn, err := grpc.Dial(target, opts...)
		if err != nil {
			log.Fatal(err)
		}
		return conn
	}}}
}

func (rpc *rpcClient) Get() *grpc.ClientConn {
	conn := rpc.pool.Get().(*grpc.ClientConn)
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = rpc.pool.New().(*grpc.ClientConn)
	}
	return conn
}
func (rpc *rpcClient) Put(conn *grpc.ClientConn) {
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		return
	}
	rpc.pool.Put(conn)
}
