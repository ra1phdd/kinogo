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

	err = application.Run()
	if err != nil {
		log.Fatal(err)
	}
}
