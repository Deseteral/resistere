package main

import (
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/webapp"
	"log"
)

func startApplication() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}

	controller.StartController(&config.Controller)

	err = webapp.StartWebServerBlocking(config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := startApplication()
	if err != nil {
		log.Println("Could not start application.")
		log.Fatal(err)
		return
	}
}
