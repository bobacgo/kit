package server

import (
	"github.com/bobacgo/kit/app"
	v1 "github.com/bobacgo/kit/examples/api/pb/user/v1"
	"github.com/bobacgo/kit/examples/internal/app/admin/service"
	"google.golang.org/grpc"
)

func GrpcRegisterServer(srv *grpc.Server, comps *app.AppOptions) {
	// register
	v1.RegisterUserServiceServer(srv, service.NewUserService())
}
