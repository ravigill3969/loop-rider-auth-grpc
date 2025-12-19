package server

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "ravigill/rider-grpc-server/proto"
)

var paymentClient pb.PaymentServiceClient

func NewPaymentClient() error {
	conn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to payment service: %w", err)
	}

	paymentClient = pb.NewPaymentServiceClient(conn)
	return nil
}

func CallCreateCheckoutSession(
	riderID string,
	estimatedPrice float32,
	pickup string,
	dropoff string,
	distance float32,
	duration int64,
	pickLat float64,
	pickLng float64,
	dropLat float64,
	dropLng float64,
	riderName string,
	riderAge int32,
	gender string,
) (*pb.CreateCheckOutSessionResponse, error) {

	if paymentClient == nil {
		return nil, fmt.Errorf("payment client not initialized (run NewPaymentClient)")
	}
	
	

	md := metadata.New(map[string]string{
		"x-internal-token": os.Getenv("INTERNAL_SERVICE_KEY"),
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	req := &pb.CreateCheckOutSessionRequest{
		RiderId:              riderID,
		RiderName:            riderName,
		RiderAge:             riderAge,
		Gender:               gender,
		EstimatedPrice:       estimatedPrice,
		PickupLocation:       pickup,
		DropoffLocation:      dropoff,
		EstimatedDistanceKm:  distance,
		EstimatedDurationMin: duration,
		PickupCoordsLatLng: &pb.Coordinates{
			Lat: pickLat,
			Lng: pickLng,
		},
		DropoffCoordsLatLng: &pb.Coordinates{
			Lat: dropLat,
			Lng: dropLng,
		},
	}

	res, err := paymentClient.CreateCheckOutSession(ctx, req)

	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("payment service error: %w", err)
	}

	return res, nil
}
