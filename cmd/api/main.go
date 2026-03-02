package main

import (
	"log"
	"os"

	"github.com/manish-npx/todo-go-echo/internal/app"
)

func main() {
	// Keep config source flexible for local/devops/docker environments.
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	apiApp, err := app.New(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer apiApp.Close()

	apiApp.Run()

}
