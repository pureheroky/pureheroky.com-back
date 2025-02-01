package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mymmrac/telego"
	"github.com/wneessen/go-mail"
	"pureheroky.com/backend/controllers"
)

func SetupRoutes(app *fiber.App, mailClient *mail.Client, bot *telego.Bot) {
	userController := controllers.ControllerService(mailClient, bot)

	app.Get("/user", userController.GetUserData)
	app.Post("/user", userController.SetUserData)

	app.Get("/skills", userController.GetSkills)
	app.Post("/skills", userController.AddSkill)

	app.Get("/projects", userController.GetProjects)
	app.Post("/projects", userController.CreateProject)

	app.Get("/commits", userController.GetLatestCommits)
	app.Post("/request", userController.SendRequest)
}
