package main

import (
	"log"
	"temporal-getting-started/internal/di"
)

func main() {
	worker, err := di.InitializeTemporal()
	if err != nil {
		log.Fatalln("cannot start worker", err)
	}
	worker.Start()

	server, err := di.InitializeApi()
	if err != nil {
		log.Fatalln("cannot start server", err)
	}
	server.Start()
}
