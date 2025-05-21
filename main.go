package main

import (
	"log"
	"temporal-getting-started/internal/di"
)

func main() {
	app, err := di.InitializeApplication()
	if err != nil {
		log.Fatalln("cannot start worker", err)
	}

	// Run worker non-blocking
	go func() {
		if err := app.Worker.Run(); err != nil {
			log.Fatalln("cannot start worker", err)
		}
	}()

	// Run server blocking
	if err = app.Server.Run(); err != nil {
		log.Fatalln("cannot start server", err)
	}
}
