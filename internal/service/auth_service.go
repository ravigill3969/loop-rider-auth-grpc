package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"ravigill/rider-grpc-server/internal/models"
	"ravigill/rider-grpc-server/internal/repository"
	"ravigill/rider-grpc-server/internal/utils"
	pb "ravigill/rider-grpc-server/proto"

	jwtlib "github.com/loop/backend/rider-auth/lib/jwt"
	"github.com/loop/backend/rider-auth/lib/middleware"

	"github.com/google/uuid"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo *repository.UserRepo
}

func NewAuthServer(repo *repository.UserRepo) *AuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	if req.User == nil {
		return &pb.AuthResponse{Success: false, Message: "User details required", Status: 400}, nil
	}

	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.User.Email)
	if existingUser != nil {
		return &pb.AuthResponse{Success: false, Message: "User already exists", Status: 409}, nil
	}

	hashedPassword, err := utils.HashPassword(req.User.Password)
	if err != nil {
		return &pb.AuthResponse{Success: false, Message: "Failed to hash password", Status: 500}, nil
	}

	user := &models.User{
		ID:          uuid.New().String(),
		Email:       req.User.Email,
		FullName:    req.User.FullName,
		Password:    hashedPassword,
		PhoneNumber: req.User.PhoneNumber,
		BirthMonth:  req.User.BirthMonth,
		BirthYear:   req.User.BirthYear,
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return &pb.AuthResponse{Success: false, Message: "Failed to create user", Status: 500}, nil
	}

	secretKeyForAccessToken := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secretKeyForAccessToken == "" {
		return &pb.AuthResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	secretKeyForRefreshToken := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if secretKeyForRefreshToken == "" {
		return &pb.AuthResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	accessToken, err := jwtlib.GenerateToken(createdUser.Email, createdUser.ID, secretKeyForAccessToken, time.Hour*24*3) // 3 days
	if err != nil {
		return &pb.AuthResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	refreshToken, err := jwtlib.GenerateToken(createdUser.Email, createdUser.ID, secretKeyForRefreshToken, time.Hour*24*7) // 7 days
	if err != nil {
		return &pb.AuthResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	return &pb.AuthResponse{
		Success: true,
		Message: "User registered successfully",
		Status:  201,
		User: &pb.User{
			Id:          createdUser.ID,
			Email:       createdUser.Email,
			FullName:    createdUser.FullName,
			PhoneNumber: createdUser.PhoneNumber,
			BirthMonth:  createdUser.BirthMonth,
			BirthYear:   createdUser.BirthYear,
			UpdatedAt:   createdUser.UpdatedAt,
			CreatedAt:   createdUser.CreatedAt,
		},
		Token: &pb.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	fmt.Println("Login Request:", req)
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &pb.LoginResponse{Success: false, Message: "Invalid credentials", Status: 401}, nil
	}

	if user == nil {
		return &pb.LoginResponse{Success: false, Message: "User not found", Status: 404}, nil
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return &pb.LoginResponse{Success: false, Message: "Invalid credentials", Status: 401}, nil
	}

	secretKeyForAccessToken := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secretKeyForAccessToken == "" {
		return &pb.LoginResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	secretKeyForRefreshToken := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if secretKeyForRefreshToken == "" {
		return &pb.LoginResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	accessToken, err := jwtlib.GenerateToken(user.Email, user.ID, secretKeyForAccessToken, time.Hour*24*3) // 3 days
	if err != nil {
		return &pb.LoginResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	refreshToken, err := jwtlib.GenerateToken(user.Email, user.ID, secretKeyForRefreshToken, time.Hour*24*7) // 7 days
	if err != nil {
		return &pb.LoginResponse{Success: false, Message: "Failed to generate token", Status: 500}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "Login successful",
		Status:  200,
		User: &pb.User{
			Id:          user.ID,
			Email:       user.Email,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			BirthMonth:  user.BirthMonth,
			BirthYear:   user.BirthYear,
			UpdatedAt:   user.UpdatedAt,
			CreatedAt:   user.CreatedAt,
		},
		Token: &pb.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		},
	}, nil
}

func (s *AuthService) GetRiderDetails(ctx context.Context, req *pb.GetRiderDetailsRequest) (*pb.GetRiderDetailsResponse, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return &pb.GetRiderDetailsResponse{Success: false, Message: "Unauthorized", Status: 401}, nil
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return &pb.GetRiderDetailsResponse{Success: false, Message: "User not found", Status: 404}, nil
	}

	return &pb.GetRiderDetailsResponse{
		Success: true,
		Message: "User details fetched successfully",
		Status:  200,
		User: &pb.User{
			Id:          user.ID,
			Email:       user.Email,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			BirthMonth:  user.BirthMonth,
			BirthYear:   user.BirthYear,
			UpdatedAt:   user.UpdatedAt,
			CreatedAt:   user.CreatedAt,
		},
	}, nil
}
