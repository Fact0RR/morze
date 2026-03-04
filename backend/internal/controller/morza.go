package controller

import (
	"fmt"
	"strconv"

	"github.com/Fact0RR/morze/internal/domain"
	"github.com/Fact0RR/morze/internal/service"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type MorzeController struct {
	service *service.MorzeService
	logger  *log.Logger
}

func NewChangeMorzeController(service *service.MorzeService, logger *log.Logger) *MorzeController {
	return &MorzeController{
		service: service,
		logger:  logger,
	}
}

func (c *MorzeController) RegisterRoutes(router fiber.Router, jwtMiddleware fiber.Handler) {
	router.Get("/messages", c.GetPrivateMessages)
	router.Post("/message", c.PostPrivateMessage)
	router.Put("/message/:id", c.UpdatePrivateMessage)
	router.Delete("/message/:id", c.DeletePrivateMessage)
}

func (c *MorzeController) GetPrivateMessages(ctx *fiber.Ctx) error {

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
	messages, err := c.service.GetPrivateMessages(ctx.Context(), contactID, limit, offset)
	if err != nil {
		c.logger.Error("err in query: ", err)
		ctx.Status(fiber.StatusInternalServerError)
	}

	return ctx.JSON(messages)
}

func (c *MorzeController) PostPrivateMessage(ctx *fiber.Ctx) error {
	var contact domain.Contact
	var err error

	c.logger.Debug("Запуск функции на публикацию приватных сообщений")
	if err = ctx.BodyParser(&contact); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse JSON",
			"details": err.Error(),
		})
	}

	id, err := c.service.PostPrivateMessages(ctx.Context(), contact.ContactID, contact.UserID, contact.Data, contact.Additionals)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot post message"),
		})
	}
	response := domain.MessageResponse{MessageID: id}
	return ctx.JSON(response)
}

func (c *MorzeController) UpdatePrivateMessage(ctx *fiber.Ctx) error {
	var updatedData domain.UpdateRequest
	messageIDString := ctx.Params("id")

	messageID, err := strconv.Atoi(messageIDString)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "message_id не цифра",
			"details": err.Error(),
		})
	}

	c.logger.Debug("Запуск функции на публикацию приватных сообщений")
	if err = ctx.BodyParser(&updatedData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse JSON",
			"details": err.Error(),
		})
	}

	if err := c.service.UpdatePrivateMessage(ctx.Context(), messageID, updatedData.UpdatedData); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot update message message_id:%d, updatedData:%s", messageID, updatedData.UpdatedData),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *MorzeController) DeletePrivateMessage(ctx *fiber.Ctx) error {
	messageIDString := ctx.Params("id")

	messageID, err := strconv.Atoi(messageIDString)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "message_id не цифра",
			"details": err.Error(),
		})
	}

	if err := c.service.DeletePrivateMessage(ctx.Context(), messageID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot delete message message_id:%d", messageID),
		})
	}
	
	return ctx.SendStatus(fiber.StatusOK)
}
