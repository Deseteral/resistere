package controller

import (
	"errors"
	"fmt"
	"testing"

	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/metrics"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/vehicle"
)

func Test_NoVehicleInRange(t *testing.T) {
	// given
	inverter := &mockInverter{state: pv.InverterState{PowerProduction: 10000, PowerConsumption: 5000}}
	vehicleCtrl := &mockVehicleController{chargingStates: map[string]vehicle.ChargingState{}}
	metrics := metrics.NewMetricsRegistry()
	ctrl := NewController(inverter, vehicleCtrl, baseConfig(0), metrics)

	// when
	ctrl.Tick()

	// then
	if len(vehicleCtrl.setAmpsCalls) != 0 {
		t.Errorf("Expected no SetChargingAmps calls, got %v", vehicleCtrl.setAmpsCalls)
	}

	// and
	if metrics.LatestFrame.PowerProductionWatts != 10000 {
		t.Errorf("Wrong power production in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
	if metrics.LatestFrame.PowerConsumptionWatts != 5000 {
		t.Errorf("Wrong power consumption in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
}

func Test_NoVehicleCharging(t *testing.T) {
	// given
	inverter := &mockInverter{state: pv.InverterState{PowerProduction: 10000, PowerConsumption: 5000}}
	vehicleCtrl := &mockVehicleController{
		chargingStates: map[string]vehicle.ChargingState{
			"VIN2": {Amps: 0, Power: 0},
		},
	}
	metrics := metrics.NewMetricsRegistry()
	ctrl := NewController(inverter, vehicleCtrl, baseConfig(0), metrics)

	// when
	ctrl.Tick()

	// then
	if len(vehicleCtrl.setAmpsCalls) != 0 {
		t.Errorf("Expected no SetChargingAmps calls, got %v", vehicleCtrl.setAmpsCalls)
	}

	// and
	if metrics.LatestFrame.PowerProductionWatts != 10000 {
		t.Errorf("Wrong power production in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
	if metrics.LatestFrame.PowerConsumptionWatts != 5000 {
		t.Errorf("Wrong power consumption in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
}

func Test_AmpsCalculation_SurplusAndCurrentAmps(t *testing.T) {
	tests := []struct {
		surplusPower    float64
		currentAmps     int
		expectedSetAmps int
		safetyMargin    int
	}{
		{+12 * 1000, 5, 16, 0},
		{+0 * 1000, 8, 8, 0},
		{-1 * 1000, 16, 15, 0},
		{-12 * 1000, 16, 5, 0},
		{+1 * 1000, 10, 10, 1000},
		{+0 * 1000, 10, 9, 1000},
		{-10, 10, 9, 1000},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("For surplus %f W and %d A currently set should set next amps to %d A", tt.surplusPower, tt.currentAmps, tt.expectedSetAmps), func(t *testing.T) {
			// given
			inverter := &mockInverter{state: pv.InverterState{PowerProduction: tt.surplusPower, PowerConsumption: 0}}
			vehicleCtrl := &mockVehicleController{
				chargingStates: map[string]vehicle.ChargingState{
					"VIN1": {Amps: tt.currentAmps, Power: float64(tt.currentAmps) * 230 * 3},
				},
			}
			ctrl := NewController(inverter, vehicleCtrl, baseConfig(tt.safetyMargin), metrics.NewMetricsRegistry())

			// when
			ctrl.Tick()

			// then
			got := vehicleCtrl.setAmpsCalls["VIN1"]
			if got != tt.expectedSetAmps {
				t.Errorf("Expected SetChargingAmps to %d, got %d", tt.expectedSetAmps, got)
			}
		})
	}
}

// Mocks
type mockInverter struct {
	state pv.InverterState
	err   error
}

func (m *mockInverter) ReadInverterState() (pv.InverterState, error) {
	return m.state, m.err
}

type mockVehicleController struct {
	chargingStates map[string]vehicle.ChargingState
	setAmpsCalls   map[string]int
	err            error
}

func (m *mockVehicleController) GetChargingState(v *vehicle.Vehicle) (*vehicle.ChargingState, error) {
	state, ok := m.chargingStates[v.Vin]
	if !ok {
		return &vehicle.ChargingState{}, errors.New("not in range")
	}
	return &state, nil
}

func (m *mockVehicleController) SetChargingAmps(v *vehicle.Vehicle, amps int) error {
	if m.setAmpsCalls == nil {
		m.setAmpsCalls = make(map[string]int)
	}
	m.setAmpsCalls[v.Vin] = amps
	return m.err
}

func baseConfig(safetyMargin int) *configuration.Config {
	return &configuration.Config{
		Controller: configuration.Controller{
			CycleIntervalSeconds: 10,
			SafetyMarginWatts:    safetyMargin,
			GridVoltage:          230,
		},
		Vehicles: configuration.Vehicles{
			Cars: []configuration.Vehicle{
				{Name: "Car1", Vin: "VIN1"},
				{Name: "Car2", Vin: "VIN2"},
			},
		},
	}
}
