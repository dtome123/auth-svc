package grpc

import (
	"fmt"
	"log"
	"net"

	"auth-svc/config"
	"auth-svc/internal/port/grpc/interceptor"
	"auth-svc/internal/services"
	"auth-svc/internal/types"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
	"github.com/dtome123/auth-sdk/jwt"
	"github.com/dtome123/auth-sdk/middlewares"
	exAuthPb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	authPb.UnimplementedTokenServiceServer
	authPb.UnimplementedAuthorizationServiceServer
	exAuthPb.UnimplementedAuthorizationServer

	cfg            *config.Config
	svc            *services.Service
	clients        map[string]types.ClientEntry
	serverVerifier jwt.Verifier
	serverSigner   jwt.Signer
	server         *grpc.Server
}

func NewGrpcServer(
	cfg *config.Config,
	svc *services.Service,
	clients map[string]types.ClientEntry,
	signer jwt.Signer,
	verifier jwt.Verifier,
) *GrpcServer {
	return &GrpcServer{
		cfg:            cfg,
		svc:            svc,
		clients:        clients,
		serverSigner:   signer,
		serverVerifier: verifier,
	}
}

func (s *GrpcServer) Run() {
	// Start listening on gRPC port
	addr := fmt.Sprintf(":%s", s.cfg.Server.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("❌ Failed to listen on %s: %v", addr, err)
	}

	// Initialize interceptors
	replayChecker := s.svc.GetReplayChecker()
	var oauthOpts []middlewares.OauthInterceptorOption
	if s.cfg.AuthConfig.EnableReplayCheck && replayChecker != nil {
		oauthOpts = append(oauthOpts, middlewares.WithGRPCReplayChecker(replayChecker))
	}

	oauthInterceptor := middlewares.NewOauthInterceptor(
		s.cfg.AuthConfig.Oauth.Aud,
		s.serverVerifier,
		oauthOpts...,
	)

	// Define interceptor rules per path/pattern
	pathRules := []interceptor.PathRule{
		{
			// Service-to-service authenticated routes
			Patterns: []string{
				authPb.TokenService_Sign_FullMethodName,
				authPb.TokenService_Refresh_FullMethodName,
				"/auth.v1.AuthorizationService/*",
			},
			Interceptors: []grpc.UnaryServerInterceptor{
				oauthInterceptor.VerifyInterceptor(),
			},
		},
	}

	// Create gRPC server with path-based unary interceptor
	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.NewPathBasedInterceptor(pathRules...),
		),
	)

	// Register services
	authPb.RegisterTokenServiceServer(s.server, s)
	authPb.RegisterAuthorizationServiceServer(s.server, s)
	exAuthPb.RegisterAuthorizationServer(s.server, s)

	// Enable reflection for easier debugging
	reflection.Register(s.server)

	log.Printf("🚀 gRPC server running at :%s\n", s.cfg.Server.GrpcPort)

	if err := s.server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
		log.Printf("❌ gRPC server exited with error: %v", err)
	}
}

func (s *GrpcServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
