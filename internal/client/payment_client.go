package client

import (
	payment_pb "ravigill/rider-grpc-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	Payment_client payment_pb.PaymentServiceClient
}

func NewPaymentClient(addr string) (*PaymentClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	return &PaymentClient{
		Payment_client: payment_pb.NewPaymentServiceClient(conn),
	}, nil
}
