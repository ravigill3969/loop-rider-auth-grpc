package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"ravigill/rider-grpc-server/internal/config"
	"ravigill/rider-grpc-server/internal/repository"
	"ravigill/rider-grpc-server/internal/service"
	pb "ravigill/rider-grpc-server/proto"

	middleware "ravigill/loop-auth-utils"
)

type RiderServer struct {
	pb.UnimplementedAuthServiceServer
}

func main() {
	err := config.LoadENV()

	if err != nil {
		fmt.Printf(".env not loaded: %v\n", err)
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = ":50052"
	} else if PORT[0] != ':' {
		PORT = ":" + PORT
	}

	db, err := config.Conn_db()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	authService := service.NewAuthServer(userRepo)

	secretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("ACCESS_TOKEN_SECRET_KEY environment variable is required")
	}

	// /<proto_package>.<ServiceName>/<MethodName>

	publicEndpoints := map[string]bool{
		"/rider_auth.AuthService/Register": true,
		"/rider_auth.AuthService/Login":    true,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(publicEndpoints, secretKey)),
	)

	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on %s", PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
