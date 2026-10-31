package port

import (
	"auth-svc/config"
	"auth-svc/internal/domain"
	"auth-svc/internal/port/grpc"
	"auth-svc/internal/port/rest"
	"auth-svc/internal/services"
	"auth-svc/internal/types"
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/dtome123/auth-sdk/jwt"
	"github.com/dtome123/auth-sdk/session"
	mongodb "github.com/dtome123/go-mongo-generic"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	cfg             *config.Config
	svc             *services.Service
	clientVerifiers map[string]types.ClientEntry
	signer          jwt.Signer
	verifier        jwt.Verifier
	userSigner      jwt.Signer
	userVerifier    jwt.Verifier
	sessionResource session.Repository[domain.Session]
}

func NewServer(cfg *config.Config, db *mongodb.Database, redisClient *redis.Client, clientVerifiers map[string]types.ClientEntry, sessionResource session.Repository[domain.Session]) *Server {

	var serverVerifier jwt.Verifier
	var serverSigner jwt.Signer

	switch cfg.AuthConfig.Oauth.Type {
	case types.AuthM2MTypeHMAC:
		serverSigner = jwt.NewHMACSigner([]byte(cfg.AuthConfig.Oauth.HMAC.Secret))
		serverVerifier = jwt.NewHMACVerifier([]byte(cfg.AuthConfig.Oauth.HMAC.Secret))
	case types.AuthM2MTypeRSA:
		var err error
		serverSigner, err = jwt.NewRS256SignerFromPath(cfg.AuthConfig.Oauth.RSA.PrivateKeyPath)
		if err != nil {
			panic(err)
		}

		serverVerifier, err = jwt.NewRS256VerifierFromPath(cfg.AuthConfig.Oauth.RSA.PublicKeyPath)
		if err != nil {
			panic(err)
		}
	}

	var userSigner jwt.Signer
	var userVerifier jwt.Verifier
	var err error
	userConfig := cfg.AuthConfig.Oauth.UserTokenConfig
	switch userConfig.Type {
	case types.AuthUserTypeHMAC:
		userSigner = jwt.NewHMACSigner([]byte(userConfig.HMAC.Secret))
		userVerifier = jwt.NewHMACVerifier([]byte(userConfig.HMAC.Secret))
	case types.AuthUserTypeRSA:
		userSigner, err = jwt.NewRS256SignerFromPath(userConfig.RSA.PrivateKeyPath)
		if err != nil {
			log.Fatalf("Invalid RSA private key for user token: %v", err)
		}
		userVerifier, err = jwt.NewRS256VerifierFromPath(userConfig.RSA.PublicKeyPath)
		if err != nil {
			log.Fatalf("Invalid RSA public key for user token: %v", err)
		}
	}

	return &Server{
		cfg:             cfg,
		svc:             services.NewService(cfg, db, redisClient, clientVerifiers, serverSigner, serverVerifier, userSigner, userVerifier, sessionResource),
		clientVerifiers: clientVerifiers,
		signer:          serverSigner,
		verifier:        serverVerifier,
		userSigner:      userSigner,
		userVerifier:    userVerifier,
		sessionResource: sessionResource,
	}
}

func (s *Server) Run() {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("❗ Recovered from panic: %v\n%s", r, debug.Stack())
		}
	}()

	grpcSvr := grpc.NewGrpcServer(s.cfg, s.svc, s.clientVerifiers, s.signer, s.verifier)
	restSvr := rest.NewRestServer(s.cfg, s.svc)

	// Run gRPC and HTTP in parallel
	go grpcSvr.Run()
	go restSvr.Run()

	// Wait for OS termination signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("⚠️ Received signal (%v). Shutting down servers gracefully...\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop REST server
	if err := restSvr.Stop(ctx); err != nil {
		log.Printf("❌ REST server shutdown error: %v\n", err)
	} else {
		log.Println("✅ REST server stopped cleanly")
	}

	// Stop gRPC server
	grpcSvr.Stop()
	log.Println("✅ gRPC server stopped cleanly")

	log.Println("👋 Server shutdown complete")
}
