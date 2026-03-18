package main

import (
	"log"

	"github.com/ramdhanrizki/bytecode-api/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("bootstrap worker app: %v", err)
	}

	if err := app.RunWorker(); err != nil {
		log.Fatalf("run worker: %v", err)
	}
}
