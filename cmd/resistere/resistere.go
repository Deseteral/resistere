package main

import (
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/webapp"
	"log"
)

func startApplication() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}

	inverter := pv.NewSolarmanInverter(&config.SolarmanInverter)

	c := controller.NewController(
		inverter,
		&config.Controller,
	)
	c.StartBackgroundTask()

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
