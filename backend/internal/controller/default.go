package controller

import (
	"encoding/json"
	"fmt"

	"github.com/Fact0RR/morze/internal/configs"
	"github.com/gofiber/fiber/v2"
)

type BaseController struct {
	settings *configs.Settings
}

func NewBaseController(settings *configs.Settings) *BaseController {
	return &BaseController{
		settings: settings,
	}
}

func (c *BaseController) Health(ctx *fiber.Ctx) error {
	return ctx.SendString("")
}

func (c *BaseController) Version(ctx *fiber.Ctx) error {
	settings := c.settings

	versionParams := map[string]string{
		"version":     "0.1.0",
		"environment": "dev",
		"commit_hash": settings.CommitHash,
	}
	if version := settings.Version; version != "" {
		versionParams["version"] = version
	}
	if commit_hash := settings.CommitHash; commit_hash != "" {
		versionParams["commit_hash"] = commit_hash
	}
	if environment := settings.Environment; environment != "" {
		versionParams["environment"] = environment
	}
	u, err := json.Marshal(versionParams)
	if err != nil {
		ctx.Locals("error", fmt.Errorf("can't encode data: %w", err))
	}

	return ctx.SendString(string(u))
}

func (c *BaseController) RegisterRotes(router fiber.Router) {
	router.Get("/version", c.Version)
	router.Get("/health", c.Health)
}
