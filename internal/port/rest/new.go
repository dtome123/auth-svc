package rest

import (
	"auth-svc/config"
	"auth-svc/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RestServer struct {
	cfg     *config.Config
	svc     *services.Service
	httpSvr *http.Server
}

func NewRestServer(cfg *config.Config, svc *services.Service) *RestServer {
	return &RestServer{
		cfg: cfg,
		svc: svc,
	}
}

func (s *RestServer) Run() {
	r := gin.Default()

	s.setupRoutes(r)

	addr := fmt.Sprintf(":%s", s.cfg.Server.HttpPort)
	s.httpSvr = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("🚀 REST server running at :%s\n", s.cfg.Server.HttpPort)

	if err := s.httpSvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("❌ REST server exited with error: %v", err)
	}
}

func (s *RestServer) Stop(ctx context.Context) error {
	if s.httpSvr != nil {
		return s.httpSvr.Shutdown(ctx)
	}
	return nil
}

func (s *RestServer) setupRoutes(r *gin.Engine) {
	r.POST("/auth/v1/token", s.Token)
}
