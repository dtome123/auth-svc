package port

import (
	"auth-svc/config"
	"auth-svc/internal/port/grpc"
	"auth-svc/internal/port/rest"
	"auth-svc/internal/services"
)

type Server struct {
	cfg *config.Config
	svc *services.Service
}

func NewServer(cfg *config.Config, svc *services.Service) *Server {
	return &Server{
		cfg: cfg,
		svc: svc,
	}
}

func (s *Server) Run() {

	grpcSvr := grpc.NewGrpcServer(s.cfg, s.svc)
	restSvr := rest.NewRestServer(s.cfg, s.svc)

	// Run gRPC and HTTP in parallel
	go grpcSvr.Run()
	go restSvr.Run()

	// Prevent main from exiting
	select {}
}
