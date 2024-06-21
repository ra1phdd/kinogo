package main

import (
	"kinogo/internal/pkg/app"
	"log"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err = application.RunREST()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = application.RunGRPC()
	if err != nil {
		log.Fatal(err)
	}
}
