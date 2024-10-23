package main

import (
	"log"
	// "net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/srijan023/jwtauth/database"
	"github.com/srijan023/jwtauth/models"
	"github.com/srijan023/jwtauth/routes"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, connErr := database.NewConnection(config)

	if connErr != nil {
		log.Fatal("There is an error during database connection")
	}

	migrateErr := models.MigrateUser(db)

	if migrateErr != nil {
		log.Fatal("Error duing migrating user to database")
	}

	routes.SetupAdminRoutes(app, db)
	routes.SetupAuthRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	app.Listen(port)
}

// TODO:
/*
User redis to store the user count
Create a account number entity on model
Create account number manually based on the user count and account type
The redis server must be append-only file type
*/
