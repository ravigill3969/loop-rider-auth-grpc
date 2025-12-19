package server

import (
	"context"
	"ravigill/rider-grpc-server/internal/service"
	pb "ravigill/rider-grpc-server/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	service *service.AuthService
}

func NewAuthServer(svc *service.AuthService) *AuthServer {
	return &AuthServer{service: svc}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return s.service.Register(ctx, req)
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.service.Login(ctx, req)
}

func (s *AuthServer) GetRiderDetails(ctx context.Context, req *pb.GetRiderDetailsRequest) (*pb.GetRiderDetailsResponse, error) {
	return s.service.GetRiderDetails(ctx, req)
}
