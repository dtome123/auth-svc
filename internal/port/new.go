package port

import (
	"auth-svc/config"
	"auth-svc/internal/port/grpc"
	"auth-svc/internal/port/rest"
	"auth-svc/internal/services"
	"auth-svc/internal/types"
	"log"
	"runtime/debug"

	"github.com/dtome123/auth-sdk/jwtutils"
	mongodb "github.com/dtome123/go-mongo-generic"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	cfg             *config.Config
	svc             *services.Service
	clientVerifiers map[string]types.ClientEntry
	signer          jwtutils.Signer
	verifier        jwtutils.Verifier
}

func NewServer(cfg *config.Config, db *mongodb.Database, redisClient *redis.Client, clientVerifiers map[string]types.ClientEntry) *Server {

	var serverVerifier jwtutils.Verifier
	var serverSigner jwtutils.Signer

	switch cfg.AuthConfig.UserJWT.Type {
	case types.AuthUserTypeHMAC:
		serverSigner = jwtutils.NewHMACSigner([]byte(cfg.AuthConfig.UserJWT.HMAC.Secret))
		serverVerifier = jwtutils.NewHMACVerifier([]byte(cfg.AuthConfig.UserJWT.HMAC.Secret))
	case types.AuthUserTypeRSA:
		var err error
		serverSigner, err = jwtutils.NewRS256SignerFromPath(cfg.AuthConfig.UserJWT.RSA.PrivateKeyPath)
		if err != nil {
			panic(err)
		}

		serverVerifier, err = jwtutils.NewRS256VerifierFromPath(cfg.AuthConfig.UserJWT.RSA.PublicKeyPath)
		if err != nil {
			panic(err)
		}
	}

	return &Server{
		cfg:             cfg,
		svc:             services.NewService(cfg, db, redisClient, clientVerifiers, serverSigner, serverVerifier),
		clientVerifiers: clientVerifiers,
		signer:          serverSigner,
		verifier:        serverVerifier,
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

	// Prevent main from exiting
	select {}
}
