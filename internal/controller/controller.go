package controller

import (
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/evse"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/vehicle"
	"log"
	"time"
)

type Controller struct {
	Vehicles []vehicle.Vehicle

	mode              Mode
	updateInterval    time.Duration
	inverter          pv.Inverter
	vehicleController vehicle.Controller
	evse              evse.Evse
}

type Mode int

const (
	ModePVAutomatic Mode = iota
	ModeManual
)

var modeName = map[Mode]string{
	ModePVAutomatic: "PV Automatic",
	ModeManual:      "Manual",
}

func (c *Controller) StartBackgroundTask() {
	log.Printf("Starting controller with %v interval.\n", c.updateInterval)

	ticker := time.NewTicker(c.updateInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.tick()
			}
		}
	}()
}

func (c *Controller) tick() {
	// The controller should only perform actions when the device is set to automatic mode.
	// It must not interfere with charging process when manual mode is set.
	if c.mode == ModeManual {
		return
	}

	// Check with EVSE if there is a car charing.
	// TODO

	log.Println("Entering controller tick.")

	// Save which car is charging and what's its current set amps.
	// If controller cannot communicate with any car (because none is in range, there was a communication error, etc.)
	// it should stop further processing.
	var selectedVehicle *vehicle.Vehicle = nil
	var currentChargingAmps int
	for _, v := range c.Vehicles {
		chargingAmps, err := c.vehicleController.GetChargingAmps(&v)
		if err != nil {
			log.Printf("Could not communicate with the car %s: %v\n.", v.Name, err)
			// Don't break the loop here. This car is probably just out of range. We should process next configured car.
		}

		if chargingAmps > 0 {
			selectedVehicle = &v
			currentChargingAmps = chargingAmps
			// We only operate on one car, so if this one is in range and charging we can skip checking other cars.
			break
		}
	}

	if selectedVehicle == nil {
		log.Println("No vehicle is charging. Exiting controller tick.")
		return
	}

	log.Printf("Selected vehicle %s with %dA set.\n", selectedVehicle.Name, currentChargingAmps)

	// Get energy surplus (kW) from inverter, change it to watts (* 1000).
	//
	// Calculate by how much we should change the charging speed.
	// Save whether surplus is positive (increase speed), or negative (decrease speed) (signum of surplus).
	// Convert surplus to its absolute value, and floor it.
	// Calculate amps:
	//   V * A = W
	//   A = W / V
	//   ^   ^   ^--this is always 230V in Europe, but it should be configurable.
	//   |   |__the energy surplus, in watts, floored.
	//   |__the amount by which we can increase charging speed.
	// TODO

	// We should call vehicle controller with
	//   saved vehicle amps +/- calculated amp diff (whether we're increasing or decreasing speed),
	//   min 5A ... max 16A.
	// TODO

	log.Println("Controlled tick finished.")
}

func (c *Controller) ChangeMode(mode Mode) {
	log.Printf("Setting controller mode to %v.\n", modeName[mode])
	c.mode = mode
}

func NewController(
	inverter pv.Inverter,
	vehicleController vehicle.Controller,
	evse evse.Evse,
	config *configuration.Config,
) Controller {
	var v []vehicle.Vehicle
	for _, c := range config.Vehicles.Cars {
		v = append(v, vehicle.Vehicle{Name: c.Name, Vin: c.Vin})
	}

	return Controller{
		Vehicles: v,

		mode:              ModeManual,
		updateInterval:    time.Duration(config.Controller.CycleIntervalSeconds) * time.Second,
		inverter:          inverter,
		vehicleController: vehicleController,
		evse:              evse,
	}
}
