package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

// RPC server
type RegisterRpcFn func(server *grpc.Server)

func RunRpcServer(addr string, register RegisterRpcFn) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	register(server)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}

// RPC Dial

var rpcClientMap = make(map[string]*grpc.ClientConn)

func RpcDial(serverName string) (*grpc.ClientConn, error) {
	if cc, ok := rpcClientMap[serverName]; ok {
		state := cc.GetState()
		if state == connectivity.Ready {
			return cc, nil
		}
	}

	// conn, err := grpc.Dial(serverName, grpc.WithInsecure())
	conn, err := grpc.Dial(serverName, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	rpcClientMap[serverName] = conn
	return conn, nil
}