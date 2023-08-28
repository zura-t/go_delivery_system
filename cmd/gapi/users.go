package gapi

import (
	"context"

	"github.com/zura-t/go_delivery_system/pb"
	"github.com/zura-t/go_delivery_system/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UsersService: %s", err)
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := c.CreateUser(ctx, &pb.CreateUserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Name:     req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateFullName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

func (server *Server) GetProfile(ctx context.Context, req *pb.UserId) (*pb.User, error) {
	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UsersService: %s", err)
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	arg := &pb.UserId{
		Id: req.GetId(),
	}
	user, err := c.GetUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UsersService: %s", err)
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	arg := &pb.UpdateUserRequest{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	user, err := c.UpdateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateFullName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}
	return violations
}

func (server *Server) AddPhone(ctx context.Context, req *pb.AddPhoneRequest) (*emptypb.Empty, error) {
	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UsersService: %s", err)
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	arg := &pb.AddPhoneRequest{
		Id:    req.GetId(),
		Phone: req.GetPhone(),
	}

	_, err = c.AddPhone(ctx, arg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (server *Server) DeleteUser(ctx context.Context, req *pb.UserId) (*emptypb.Empty, error) {
	conn, err := grpc.Dial(server.config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UsersService: %s", err)
	}
	defer conn.Close()

	c := pb.NewUsersServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	arg := &pb.UserId{
		Id: req.GetId(),
	}

	_, err = c.DeleteUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
