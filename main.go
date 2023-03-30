package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shahriarsohan/new_blog/initializers"
	"github.com/shahriarsohan/new_blog/routes"
)

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	initializers.ConnectDB(&config)

	if err != nil {
		log.Fatal("Error loading env file")
	}

	app := fiber.New()
	routes.Setup(app)
	app.Listen(":" + config.ServerPort)
}
