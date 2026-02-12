package client

import (
	trippb "ravigill/loop-grpc-trip/proto/trip_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TripClient struct {
	trip_client trippb.TripServiceClient
}

func NewTripClient(addr string) (*TripClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	
	if err != nil {
		return nil, err
	}

	return &TripClient{
		trip_client: trippb.NewTripServiceClient(conn),
	}, nil
}

