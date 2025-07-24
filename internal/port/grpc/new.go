package grpc

import (
	"auth-svc/config"
	"auth-svc/internal/services"
	"log"
	"net"

	"google.golang.org/grpc"

	authPb "github.com/dtome123/auth-sdk/gen/go/auth/v1"
	exAuthPb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type GrpcServer struct {
	authPb.UnimplementedAuthServiceServer
	exAuthPb.UnimplementedAuthorizationServer

	cfg *config.Config
	svc *services.Service
}

func NewGrpcServer(cfg *config.Config, svc *services.Service) *GrpcServer {
	return &GrpcServer{
		cfg: cfg,
		svc: svc,
	}
}

func (s *GrpcServer) Run() {
	lis, err := net.Listen("tcp:", s.cfg.Server.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	authPb.RegisterAuthServiceServer(grpcServer, s)
	exAuthPb.RegisterAuthorizationServer(grpcServer, s)

	log.Println("🚀 gRPC server running at :", s.cfg.Server.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
