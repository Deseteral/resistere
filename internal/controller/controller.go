package controller

import (
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/vehicle"
	"log"
	"time"
)

type Controller struct {
	updateInterval    time.Duration
	inverter          pv.Inverter
	vehicleController vehicle.Controller
}

func (c *Controller) StartBackgroundTask() {
	log.Printf("Starting controller with %v interval.\n", c.updateInterval)

	// Run first tick before ticker starts.
	c.tick()

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
	log.Println("Starting controller tick.")

	// Check which car is in range and is charging.
	// Save which one is charging and what is its current set amps.
	// If none, log and return.

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

	// We should call vehicle controller with
	//   saved vehicle amps +/- calculated amp diff (whether we're increasing or decreasing speed),
	//   min 5A ... max 16A.
}

func NewController(
	inverter pv.Inverter,
	vehicleController vehicle.Controller,
	config *configuration.Controller,
) Controller {
	return Controller{
		updateInterval:    time.Duration(config.CycleIntervalSeconds) * time.Second,
		inverter:          inverter,
		vehicleController: vehicleController,
	}
}
