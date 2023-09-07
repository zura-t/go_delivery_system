package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/zura-t/go_delivery_system/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
)

var ignoreMethods = []string{pb.UsersService_CreateUser_FullMethodName, pb.UsersService_LoginUser_FullMethodName}

func (server *Server) AuthorizeUser(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		for _, imethod := range ignoreMethods {
			if info.FullMethod == imethod {
				md, ok := metadata.FromIncomingContext(ctx)
				if !ok {
					return nil, fmt.Errorf("missing metadata")
				}

				values := md.Get(authorizationHeader)
				if len(values) == 0 {
					return nil, fmt.Errorf("missing authorization header")
				}

				authHeader := values[0]
				fields := strings.Fields(authHeader)
				if len(fields) < 2 {
					return nil, fmt.Errorf("invalid authorization header")
				}

				authType := strings.ToLower(fields[0])
				if authType != authorizationTypeBearer {
					return nil, fmt.Errorf("invalid authorization type: %s", authType)
				}

				accessToken := fields[1]
				payload, err := server.tokenMaker.VerifyToken(accessToken)
				if err != nil {
					return nil, fmt.Errorf("invalid access token: %s", err)
				}

				return payload, nil
			}
		}
		return nil, nil
}
