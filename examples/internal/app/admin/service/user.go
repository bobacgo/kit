package service

import (
	"context"

	v1 "github.com/bobacgo/kit/examples/api/pb/user/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	v1.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u UserService) GetUserById(ctx context.Context, request *v1.GetUserRequest) (*v1.UserResponse, error) {
	_, span := otel.Tracer("examples-service").Start(ctx, "GetUserById")
	defer span.End()

	user := &v1.UserResponse{
		User: &v1.User{
			Id:       request.Id,
			Username: "test",
		},
	}

	return user, nil
}

func (u UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.UserResponse, error) {
	// if err := validator.StructCtx(ctx, request); err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, err.Error())
	// }
	return nil, nil
}

func (u UserService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (u UserService) DeleteUser(ctx context.Context, request *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	// TODO implement me
	return &emptypb.Empty{}, nil
}
