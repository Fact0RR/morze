package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

var (
	errInvalidJWTToken   = errors.New("invalid JWT Token")
	errAuthHeaderIsEmpty = errors.New("\"Authorization\" header is empty")
)

func AuthMiddleware(secretKey []byte, logger *log.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if err := parseJWT(secretKey, authHeader, logger); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Next()
	}
}

func parseJWT(secretKey []byte, authHeader string, logger *log.Logger) error {
	if len(authHeader) > 0 {
		authToken := strings.Replace(authHeader, "Bearer ", "", 1)
		tokenJWT, err := jwt.Parse(authToken, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Warnf("Отказ в доступе: неподдерживаемый метод подписи токена %s:%s", "alg", token.Header["alg"].(string))
				return nil, errInvalidJWTToken
			}
			return secretKey, nil
		})
		if err != nil {
			logger.Errorf("Error parse jwt token  error %s ", err)
			return errInvalidJWTToken
		}

		if !tokenJWT.Valid {
			logger.Warn("Отказ в доступе: невалидный токен")
			return errInvalidJWTToken
		}

		return nil
	} else {
		logger.Debug("Error parse jwt token - Auth header is empty")

		return errAuthHeaderIsEmpty
	}
}
