package handlers

import (
	"context"
	"log"

	"github.com/Tel3scop/auth/internal/services/user_service"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userServer struct {
	userAPI.UnimplementedUserV1Server
}

func (s *userServer) Get(ctx context.Context, req *userAPI.GetRequest) (*userAPI.GetResponse, error) {
	log.Printf("Getting user id: %d", req.GetId())
	user, err := user_service.GetByID(ctx, req.Id)
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
	}, err
}

func (s *userServer) Create(ctx context.Context, req *userAPI.CreateRequest) (*userAPI.CreateResponse, error) {
	log.Printf("Creating data: %+v", req)
	createdUserID, err := user_service.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &userAPI.CreateResponse{
		Id: createdUserID,
	}, err
}

func (s *userServer) Update(ctx context.Context, req *userAPI.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Updating data: %+v", req)
	err := user_service.UpdateByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}

func (s *userServer) Delete(ctx context.Context, req *userAPI.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Deleting data: %+v", req)
	err := user_service.DeleteByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}
