package main

import (
	"log"

	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/evse"
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

	var inverterImpl pv.Inverter
	var vehicleControllerImpl vehicle.Controller
	var evseImpl evse.Evse

	if config.SimulatorMode {
		log.Println("Running in simulator mode.")
		inverterImpl = pv.NewSimulatedInverter()
		vehicleControllerImpl = vehicle.NewSimulatedVehicleController()
		evseImpl = evse.NewSimulatedEvse()
	} else {
		inverterImpl = pv.NewSolarmanInverter(&config.SolarmanInverter)
		vehicleControllerImpl = vehicle.NewTeslaControlController(&config.TeslaControl)
		evseImpl = evse.NewTeslaWallConnector(&config.TeslaWallConnector)
	}

	mr := metrics.NewMetricsRegistry()

	c := controller.NewController(
		inverterImpl,
		vehicleControllerImpl,
		evseImpl,
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
