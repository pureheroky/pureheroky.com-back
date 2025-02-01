package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"github.com/wneessen/go-mail"
	"pureheroky.com/backend/database"
	"pureheroky.com/backend/routes"
)

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mailClient, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(os.Getenv("EMAIL_FROM")),
		mail.WithPassword(os.Getenv("EMAIL_PASS")),
	)

	if err != nil {
		log.Fatalf("Failed to create mail client: %s", err)
	}

	bot, err := telego.NewBot(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create new bot: %s", err)
	}

	err = database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}

	app.Static("/images", "../static/images")

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://pureheroky.com",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	routes.SetupRoutes(app, mailClient, bot)

	app.Listen("127.0.0.1:8080")
}
