package main

import (
	"blog_starter_project_g66/Repositories"
	"log"
)

func main() {
	mongoClient, err := repositories.Connect()
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer mongoClient.Disconnect()
}
