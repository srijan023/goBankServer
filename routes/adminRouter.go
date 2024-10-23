package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/srijan023/jwtauth/controllers"
	"github.com/srijan023/jwtauth/middleware"
	"gorm.io/gorm"
)

func SetupAdminRoutes(app *fiber.App, db *gorm.DB) {
	// api := app.Group("/user", middleware.Authenticate())
	api := app.Group("/admin")
	ur := controllers.UserRouter{
		DB: db,
	}

	api.Use("/", middleware.Authenticate())
	api.Get("/user/:id", ur.GetUser())
	api.Get("/user", ur.GetUsers())
}
