package grpc

import (
	"fmt"
	"log"
	"net"
	"time"

	"auth-svc/config"
	"auth-svc/internal/port/grpc/interceptor"
	"auth-svc/internal/services"
	"auth-svc/internal/types"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
	"github.com/dtome123/auth-sdk/jwtutils"
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
	serverVerifier jwtutils.Verifier
	serverSigner   jwtutils.Signer
}

func NewGrpcServer(
	cfg *config.Config,
	svc *services.Service,
	clients map[string]types.ClientEntry,
	signer jwtutils.Signer,
	verifier jwtutils.Verifier,
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
	internalInterceptor := interceptor.NewInternalDirectInterceptor(s.cfg.AuthConfig, s.clients)
	oauthOpts := []middlewares.OauthInterceptorOption{
		middlewares.WithMethodRules(getAuthorizationMethodRules()),
	}
	if s.cfg.AuthConfig.EnableReplayCheck {
		oauthOpts = append(oauthOpts, middlewares.WithGRPCReplayChecker(5*time.Minute))
	}

	oauthInterceptor := middlewares.NewOauthInterceptor(
		s.cfg.AuthConfig.Aud,
		s.serverVerifier,
		oauthOpts...,
	)

	// Create gRPC server with chained unary interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			internalInterceptor.UnaryInterceptor(),
			oauthInterceptor.UnaryInterceptor(),
		),
	)

	// Register services
	authPb.RegisterTokenServiceServer(grpcServer, s)
	authPb.RegisterAuthorizationServiceServer(grpcServer, s)
	exAuthPb.RegisterAuthorizationServer(grpcServer, s)

	// Enable reflection for easier debugging
	reflection.Register(grpcServer)

	log.Printf("🚀 gRPC server running at :%s\n", s.cfg.Server.GrpcPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve gRPC: %v", err)
	}
}

// getAuthorizationMethodRules returns all method rules that require service-level auth.
func getAuthorizationMethodRules() []middlewares.MethodRule {
	methods := []string{
		authPb.AuthorizationService_GetUserPermissions_FullMethodName,
		authPb.AuthorizationService_HasPermission_FullMethodName,
		authPb.AuthorizationService_CheckPermissions_FullMethodName,
		authPb.AuthorizationService_ListPermissionPaths_FullMethodName,
		authPb.AuthorizationService_AssignRolesToUser_FullMethodName,
		authPb.AuthorizationService_GetUserRoles_FullMethodName,
		authPb.AuthorizationService_CreateRole_FullMethodName,
		authPb.AuthorizationService_UpdateRole_FullMethodName,
		authPb.AuthorizationService_ListRole_FullMethodName,
	}

	rules := make([]middlewares.MethodRule, len(methods))
	for i, method := range methods {
		rules[i] = middlewares.MethodRule{
			Pattern: method,
			Mode:    middlewares.AuthService,
		}
	}
	return rules
}
