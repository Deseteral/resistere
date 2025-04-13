package main

import (
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/evse"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/vehicle"
	"github.com/deseteral/resistere/internal/webapp"
	"log"
)

func startApplication() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}

	var inverter pv.Inverter
	var vehicleController vehicle.Controller
	var wallbox evse.Evse

	if config.SimulatorMode {
		log.Println("Running in simulator mode.")
		inverter = pv.NewSimulatedInverter()
		vehicleController = vehicle.NewSimulatedVehicleController()
		wallbox = evse.NewSimulatedEvse()
	} else {
		inverter = pv.NewSolarmanInverter(&config.SolarmanInverter)
		vehicleController = vehicle.NewTeslaControlController(&config.TeslaControl)
		wallbox = evse.NewTeslaWallConnector(&config.TeslaWallConnector)
	}

	c := controller.NewController(
		inverter,
		vehicleController,
		wallbox,
		&config.Controller,
	)
	c.ChangeMode(controller.ModePVAutomatic) // TODO: This should be controlled by physical toggle switch.
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
