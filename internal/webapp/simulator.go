package webapp

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/evse"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/vehicle"
	"github.com/deseteral/resistere/internal/webapp/view"
)

// Simulator holds the same simulated devices used by the charging controller.
// It is nil when the application is connected to real hardware.
type Simulator struct {
	Inverter *pv.SimulatedInverter
	EVSE     *evse.SimulatedEvse
	Vehicles *vehicle.SimulatedVehicleController
}

func (s *Simulator) handler(config *configuration.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := ""
		if r.Method == http.MethodPost {
			if err := s.update(r, config.Vehicles.Cars); err != nil {
				message = err.Error()
			} else {
				message = "Simulator parameters applied. The controller reads them on its next scheduled tick."
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Header.Get("HX-Request") == "true" {
			view.SimulatorControls(s.Inverter, s.EVSE, s.Vehicles, config.Vehicles.Cars, message).Render(r.Context(), w)
		} else {
			view.SimulatorPage(s.Inverter, s.EVSE, s.Vehicles, config.Vehicles.Cars, message).Render(r.Context(), w)
		}
	}
}

func (s *Simulator) update(r *http.Request, cars []configuration.Vehicle) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Invalid form data. No parameters changed.")
	}
	production, err := strconv.ParseFloat(r.PostForm.Get("production"), 64)
	if err != nil || math.IsNaN(production) || math.IsInf(production, 0) || production < 0 {
		return fmt.Errorf("Production must be a finite, non-negative number of watts. No parameters changed.")
	}
	consumption, err := strconv.ParseFloat(r.PostForm.Get("consumption"), 64)
	if err != nil || math.IsNaN(consumption) || math.IsInf(consumption, 0) || consumption < 0 {
		return fmt.Errorf("Consumption must be a finite, non-negative number of watts. No parameters changed.")
	}
	states := make([]vehicle.SimulatedVehicleState, len(cars))
	for i, car := range cars {
		prefix := fmt.Sprintf("vehicle_%d_", i)
		amps, err := strconv.Atoi(r.PostForm.Get(prefix + "amps"))
		if err != nil || amps < 0 || amps > 16 {
			return fmt.Errorf("%s: current must be an integer from 0 to 16 A. No parameters changed.", car.Name)
		}
		power, err := strconv.Atoi(r.PostForm.Get(prefix + "power"))
		if err != nil || power < 0 || power > 11 {
			return fmt.Errorf("%s: reported power must be an integer from 0 to 11 kW. No parameters changed.", car.Name)
		}
		states[i] = vehicle.SimulatedVehicleState{
			Reachable:     r.PostForm.Get(prefix+"reachable") == "on",
			ChargingState: vehicle.ChargingState{Amps: amps, Power: power},
		}
	}
	// Validate the entire form before updating any of the devices.
	s.Inverter.SetState(pv.InverterState{PowerProduction: production, PowerConsumption: consumption})
	s.EVSE.SetConnected(r.PostForm.Get("connected") == "on")
	for i, car := range cars {
		s.Vehicles.SetSimulationState(car.Vin, states[i])
	}
	return nil
}
