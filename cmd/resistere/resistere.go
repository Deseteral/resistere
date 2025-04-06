package main

import (
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/inverter"
	"github.com/deseteral/resistere/internal/webapp"
	"log"
)

func startApplication() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}

	c := controller.NewController(
		inverter.NewSolarmanInverter(&config.Inverter),
		&config.Controller,
	)
	c.StartController()

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
