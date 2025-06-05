package controller

import (
	"log"
	"math"
	"time"

	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/evse"
	"github.com/deseteral/resistere/internal/metrics"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/utils"
	"github.com/deseteral/resistere/internal/vehicle"
)

type Controller struct {
	Vehicles []vehicle.Vehicle
	Mode     Mode

	updateInterval    time.Duration
	inverter          pv.Inverter
	vehicleController vehicle.Controller
	evse              evse.Evse
	metricsRegistry   *metrics.Registry
	config            *configuration.Config
}

func (c *Controller) StartBackgroundTask() {
	log.Printf("Starting controller with %v interval.\n", c.updateInterval)

	// TODO: Currently every tick starts at set interval. This is not correct behaviour.
	//       What should actually happen is each new tick should start at X second interval after the previous one ended
	//       and not when it started.
	//       If processing tick is longer then the interval, then the next tick will start immediately.
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

type processingData struct {
	metricsFrame metrics.Frame

	selectedVehicle             *vehicle.Vehicle
	selectedVehicleChargingAmps int
	selectedVehicleMetricsFrame *metrics.VehicleFrame

	energySurplus float64
	nextAmps      int
}

func (c *Controller) tick() {
	log.Println("Entering controller tick.")

	data := processingData{}

	// Setup metrics frame for current processing tick.
	// Make sure that all exit points of the tick commit the new metrics frame as the latest one.
	data.metricsFrame = metrics.Frame{
		Timestamp: time.Now(),
	}
	defer func() { c.metricsRegistry.LatestFrame = data.metricsFrame }()

	// Read and calculate the current energy surplus from the inverter.
	//
	// This is done before checking the controller mode, because we want to have current production
	// and consimption for metrics regardles of the controller mode and charging state.
	energySurplusErr := c.calculateEnergySurplus(&data)

	// The controller should only perform actions when the device is set to automatic mode.
	// It must not interfere with charging process when manual mode is set.
	if c.Mode == ModeManual {
		log.Println("Exiting controller tick because it's in manual mode.")
		return
	}

	// Check for all known vehicles and see which one (if any) is charging.
	c.selectVehicleForProcessing(&data)

	// If controller cannot communicate with any car (because none is in range, there was a communication error, etc.)
	// it should stop further processing.
	if data.selectedVehicle == nil {
		log.Println("No vehicle is charging. Exiting controller tick.")
		return
	}
	log.Printf("Selected vehicle %s with %dA set.\n", data.selectedVehicle.Name, data.selectedVehicleChargingAmps)

	// If controller cannot get energy surplus data from intverter it should stop further processing.
	if energySurplusErr != nil {
		log.Printf("Could not read energy surplus from inverter: %v. Exiting controller tick.\n", energySurplusErr)
		return
	}

	// Perform the amperage calculation and send the derived value to the selected vehicle.
	c.calculateNextAmpsToBeSet(&data)
	c.setChargingAmpsToSelectedVehicle(&data)

	log.Println("Controlled tick finished successfully.")
}

func (c *Controller) selectVehicleForProcessing(data *processingData) {
	// Save which car is charging and what's its current set amps.

	// TODO: Perhaps this could be running in parallel.

	for _, v := range c.Vehicles {
		vehicleMetricsFrame := metrics.VehicleFrame{
			Name: v.Name,
		}

		chargingState, err := c.vehicleController.GetChargingState(&v)

		if err != nil {
			log.Printf("Could not communicate with the car %s: %v.\n", v.Name, err)

			// Don't break the loop here.
			// This car is probably just out of range. We should process next configured car.
			continue
		}

		vehicleMetricsFrame.IsInRange = true

		if chargingState.Amps > 0 {
			vehicleMetricsFrame.SetChargingAmps = chargingState.Amps
			vehicleMetricsFrame.ChargingPowerWatts = chargingState.Power * 1000.0

			data.selectedVehicle = &v
			data.selectedVehicleChargingAmps = chargingState.Amps
			data.selectedVehicleMetricsFrame = &vehicleMetricsFrame
		} else {
			log.Printf("Car %s is in-range but not charging.\n", v.Name)
		}

		data.metricsFrame.VehicleFrames = append(data.metricsFrame.VehicleFrames, vehicleMetricsFrame)
	}
}

func (c *Controller) calculateEnergySurplus(data *processingData) error {
	// Get energy surplus from inverter.
	inverterState, err := c.inverter.ReadEnergySurplus()
	if err != nil {
		return err
	}

	data.metricsFrame.PowerProductionWatts = inverterState.PowerProduction
	data.metricsFrame.PowerConsumptionWatts = inverterState.PowerConsumption

	data.energySurplus = math.Floor(inverterState.PowerProduction - inverterState.PowerConsumption)

	log.Printf("Current energy surplus: %d W.\n", int(data.energySurplus))

	return nil
}

func (c *Controller) calculateNextAmpsToBeSet(data *processingData) {
	// Add safety margin to surplus, to ensure that we don't charge with energy from the grid.
	data.energySurplus -= float64(c.config.Controller.SafetyMarginWatts)

	// Calculate by how much we should change the charging speed.
	//   3 * V * A = W
	//   A = W / (V * 3)
	//   ▲   ▲    ▲   ▲
	//   │   │    │   └──EVSE is using three-phases to charge.
	//   │   │    └──the electric potential of energy grid, in volts.
	//   │   └──the energy surplus (the difference between power generated by PV and total power used), in watts.
	//   └──the amount by which we can increase or decrease charging current, in amps.
	gridVoltage := 230 // TODO: This is always 230V or so in Europe, but it should be configurable.
	deltaAmps := int(data.energySurplus) / (gridVoltage * 3)

	log.Printf("Calculated delta amps: %dA.\n", deltaAmps)

	// Calculate charging amps the car should use.
	// It has to be between 5A and 16A (min and max charging amps for 3kW and 11kW).
	data.nextAmps = utils.Clamp(data.selectedVehicleChargingAmps+deltaAmps, 5, 16)
}

func (c *Controller) setChargingAmpsToSelectedVehicle(data *processingData) error {
	// Send next charging amps value to the car.
	log.Printf("Setting charging amps to %d A for car %s.\n", data.nextAmps, data.selectedVehicle.Name)

	err := c.vehicleController.SetChargingAmps(data.selectedVehicle, data.nextAmps)

	if err != nil {
		log.Printf("Could not set charging amps for car %s: %v. Exiting controller tick.\n", data.selectedVehicle.Name, err)
		return err
	}

	data.selectedVehicleMetricsFrame.SetChargingAmps = data.nextAmps

	return nil
}

// TODO: When set to manual it should return back to automatic after certain time (10 minutes?) of not charging.
//
// This is important as "automatic" is the default, and leaving the controller running in "manual"
// could lead to accidental grid usage.
func (c *Controller) ChangeMode(mode Mode) {
	log.Printf("Setting controller mode to %v.\n", modeName[mode])
	c.Mode = mode
}

func NewController(
	inverter pv.Inverter,
	vehicleController vehicle.Controller,
	config *configuration.Config,
	metricsRegistry *metrics.Registry,
) Controller {
	var v []vehicle.Vehicle
	for _, c := range config.Vehicles.Cars {
		v = append(v, vehicle.Vehicle{Name: c.Name, Vin: c.Vin})
	}

	return Controller{
		Vehicles: v,
		Mode:     ModeManual,

		updateInterval:    time.Duration(config.Controller.CycleIntervalSeconds) * time.Second,
		inverter:          inverter,
		vehicleController: vehicleController,
		metricsRegistry:   metricsRegistry,
		config:            config,
	}
}
