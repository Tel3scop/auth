package access

import (
	"context"
	"fmt"
	"strings"

	accessAPI "github.com/Tel3scop/auth/pkg/access_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check handler для создания нового пользователя
func (i *Implementation) Check(ctx context.Context, request *accessAPI.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, fmt.Errorf("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], i.cfg.Encrypt.AuthPrefix) {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], i.cfg.Encrypt.AuthPrefix)
	err := i.accessService.Check(ctx, request.EndpointAddress, accessToken)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
