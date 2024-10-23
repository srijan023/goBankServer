package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/srijan023/jwtauth/helpers"
)

func Authenticate() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		clientToken := ctx.Get("Authorization")
		if clientToken == "" {
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "You need to pass an authentication token"})
		}

		claims, err := helpers.ValidateToken(ctx, clientToken)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": err.Error()})
		}

		ctx.Set("email", claims.Email)
		ctx.Set("id", claims.Id)

		return ctx.Next()
	}

}
