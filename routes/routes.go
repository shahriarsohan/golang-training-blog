package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shahriarsohan/new_blog/controllers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api", logger.New())

	api.Post("register", controllers.SignUpUser)
	api.Post("login", controllers.Login)
	api.Get("verifyemail/:verificationCode", controllers.VerifyEmail)
}
