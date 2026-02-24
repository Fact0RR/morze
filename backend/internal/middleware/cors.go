package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",                                           // Разрешенные источники
		AllowMethods:     "GET,POST,PUT,DELETE",                         // Разрешенные методы
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization", // Разрешенные заголовки
		//AllowCredentials: true,
	})
}
