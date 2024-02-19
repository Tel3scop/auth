package handlers

import (
	"context"

	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/services"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userServer struct {
	userAPI.UnimplementedUserV1Server
	services *services.Container
	cfg      *config.Config
}

func (s *userServer) Get(ctx context.Context, req *userAPI.GetRequest) (*userAPI.GetResponse, error) {
	user, err := s.services.Users.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userAPI.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Role:      user.Role,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (s *userServer) Create(ctx context.Context, req *userAPI.CreateRequest) (*userAPI.CreateResponse, error) {
	createdUserID, err := s.services.Users.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &userAPI.CreateResponse{Id: createdUserID}, nil
}

func (s *userServer) Update(ctx context.Context, req *userAPI.UpdateRequest) (*emptypb.Empty, error) {
	err := s.services.Users.UpdateByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *userServer) Delete(ctx context.Context, req *userAPI.DeleteRequest) (*emptypb.Empty, error) {
	err := s.services.Users.DeleteByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
