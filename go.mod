module ravigill/rider-grpc-server

go 1.25.1

require ravigill/loop-auth-utils v0.0.0

replace ravigill/loop-auth-utils => ../../libs/authutils

require ravigill/loop-grpc-trip v0.0.0

replace ravigill/loop-grpc-trip => ../../trip/grpc-trip

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.45.0
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

require github.com/google/uuid v1.6.0

require (
	github.com/lib/pq v1.10.9
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
)
