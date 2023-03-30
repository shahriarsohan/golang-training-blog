package main

import (
	"fmt"
	"log"

	"github.com/shahriarsohan/new_blog/initializers"
	"github.com/shahriarsohan/new_blog/models"
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("? Could not load environment configs", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	fmt.Println("? Migration complete")
}
