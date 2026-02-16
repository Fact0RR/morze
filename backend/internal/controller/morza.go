package controller

import (
	"strconv"

	"github.com/Fact0RR/morza/internal/domain"
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

func (c *MorzaController) RegisterRoutes(router fiber.Router, jwtMiddleware fiber.Handler) {
	router.Get("/messages", c.GetPrivateMessages)
	router.Post("/message", c.PostPrivateMessage)
}

func (c *MorzaController) GetPrivateMessages(ctx *fiber.Ctx) error {

	contactID, err := strconv.Atoi(ctx.Query("contact"))
	if err != nil {
		c.logger.Error("contact is not int: ", err)
		ctx.Status(fiber.StatusBadRequest)
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		c.logger.Error("limit is not int: ", err)
		ctx.Status(fiber.StatusBadRequest)
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		c.logger.Error("offset is not int: ", err)
		ctx.Status(fiber.StatusBadRequest)
	}

	c.logger.Debug("Запуск функции на получение приватных сообщений")
	messages, err:= c.service.GetPrivateMessages(ctx.Context(), contactID, limit, offset)
	if err != nil {
		c.logger.Error("err in query: ", err)
		ctx.Status(fiber.StatusInternalServerError)
	}

	return ctx.JSON(messages)
}

func (c *MorzaController) PostPrivateMessage(ctx *fiber.Ctx) error {
	var contact domain.Contact

	c.logger.Debug("Запуск функции на публикацию приватных сообщений")
	if err := ctx.BodyParser(&contact); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Cannot parse JSON",
            "details": err.Error(),
        })
    }

	return ctx.JSON(contact)
}
