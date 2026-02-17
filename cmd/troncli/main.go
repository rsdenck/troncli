package main

import (
	"log"

	"github.com/mascli/troncli/internal/ui"
)

func main() {
	app := ui.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
