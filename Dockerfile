FROM golang:1.25.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY ./rider-auth/auth-grpc/go.mod ./rider-auth/auth-grpc/go.sum ./rider-auth/auth-grpc/

COPY libs ./libs
COPY trip ./trip

WORKDIR /app/rider-auth/auth-grpc

RUN go mod download

COPY ./rider-auth/auth-grpc .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rider-auth-grpc ./cmd/grpc-server

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/rider-auth/auth-grpc .

EXPOSE 50052
USER nonroot:nonroot

CMD ["./rider-auth-grpc"]
