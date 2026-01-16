package rest

import (
	"auth-svc/config"
	"auth-svc/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

type RestServer struct {
	cfg *config.Config
	svc *services.Service
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

	log.Println("🚀 REST server running at :", s.cfg.Server.HttpPort)

	r.Run(":8080")
}

func (s *RestServer) setupRoutes(r *gin.Engine) {
	r.POST("/auth/v1/token", s.Token)
}
