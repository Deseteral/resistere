package main

import (
	"log"

	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/metrics"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/vehicle"
	"github.com/deseteral/resistere/internal/webapp"
)

func startApplication() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}

	var inverter pv.Inverter
	var vehicleController vehicle.Controller

	if config.SimulatorMode {
		log.Println("Running in simulator mode.")
		inverter = pv.NewSimulatedInverter()
		vehicleController = vehicle.NewSimulatedVehicleController()
	} else {
		inverter = pv.NewSolarmanInverter(&config.SolarmanInverter)
		vehicleController = vehicle.NewTeslaControlController(&config.TeslaControl)
	}

	mr := metrics.NewMetricsRegistry()

	c := controller.NewController(
		inverter,
		vehicleController,
		config,
		mr,
	)

	c.ChangeMode(controller.ModePVAutomatic)

	c.StartBackgroundTask()

	err = webapp.StartWebServerBlocking(config, &c, mr)
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
