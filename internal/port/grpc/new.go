package grpc

import (
	"fmt"
	"log"
	"net"

	"auth-svc/config"
	"auth-svc/internal/port/grpc/interceptor"
	"auth-svc/internal/services"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
	exAuthPb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Server.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	clientAssertionInterceptor := interceptor.NewUserDelegationInterceptor(s.cfg.AuthConfig)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(clientAssertionInterceptor.UnaryInterceptor()),
	)

	reflection.Register(grpcServer)

	authPb.RegisterAuthServiceServer(grpcServer, s)
	exAuthPb.RegisterAuthorizationServer(grpcServer, s)

	log.Println("🚀 gRPC server running at :", s.cfg.Server.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
