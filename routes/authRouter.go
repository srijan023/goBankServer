package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/srijan023/jwtauth/controllers"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/auth")
	ur := controllers.UserRouter{
		DB: db,
	}

	api.Post("/signup", ur.Signup())
	api.Post("/login", ur.Login())
}
