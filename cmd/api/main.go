package main

import (
	"log"

	"github.com/ramdhanrizki/bytecode-api/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("bootstrap api app: %v", err)
	}

	if err := app.RunAPI(); err != nil {
		log.Fatalf("run api: %v", err)
	}
}
