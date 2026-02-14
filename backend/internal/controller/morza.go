package controller

import (
	"github.com/Fact0RR/morza/internal/service"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type MorzaController struct {
	service *service.MorzaService
	logger  *log.Logger
}

func NewChangeMorzaController(service *service.MorzaService, logger *log.Logger) *MorzaController {
	return &MorzaController{
		service: service,
		logger:  logger,
	}
}

func (c *MorzaController) RegisterRotes(router fiber.Router, jwtMiddleware fiber.Handler) {}
