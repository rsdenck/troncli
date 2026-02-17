package main

import (
	"log"

	"github.com/mascli/troncli/internal/ui"
)

func main() {
	app, err := ui.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
