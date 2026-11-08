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
	evse := &mockEvse{connected: true}
	metrics := metrics.NewMetricsRegistry()
	ctrl := NewController(inverter, vehicleCtrl, evse, baseConfig(0), metrics)

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
	evse := &mockEvse{connected: true}
	metrics := metrics.NewMetricsRegistry()
	ctrl := NewController(inverter, vehicleCtrl, evse, baseConfig(0), metrics)

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
			chargingPower := int(float64(tt.currentAmps) * 230 * 3 / 1000)
			vehicleCtrl := &mockVehicleController{
				chargingStates: map[string]vehicle.ChargingState{
					"VIN2": {Amps: tt.currentAmps, Power: chargingPower},
				},
			}
			evse := &mockEvse{connected: true}
			metrics := metrics.NewMetricsRegistry()
			ctrl := NewController(inverter, vehicleCtrl, evse, baseConfig(tt.safetyMargin), metrics)

			// when
			ctrl.Tick()

			// then
			got := vehicleCtrl.setAmpsCalls["VIN2"]
			if got != tt.expectedSetAmps {
				t.Errorf("Expected SetChargingAmps to %d, got %d", tt.expectedSetAmps, got)
			}
			if vehicleCtrl.setAmpsCalls["VIN1"] != 0 {
				t.Errorf("Not charging vehicle should not be having amps set")
			}

			// and
			metric := metrics.LatestFrame.VehicleFrames[0].ChargingPowerWatts
			expectedMetric := float64(chargingPower * 1000.0)
			if metric != expectedMetric {
				t.Errorf("Expected charging power in metrics to be %f, got %f", expectedMetric, metric)
			}
			if len(metrics.LatestFrame.VehicleFrames) != 1 {
				t.Errorf("Only charging vehicles should be included in metrics")
			}
		})
	}
}

func Test_EVSE_Error(t *testing.T) {
	// given
	inverter := &mockInverter{state: pv.InverterState{PowerProduction: 10000, PowerConsumption: 5000}}
	vehicleCtrl := &mockVehicleController{chargingStates: map[string]vehicle.ChargingState{}}
	evse := &mockEvse{connected: false, err: errors.New("evse communication error")}
	metrics := metrics.NewMetricsRegistry()
	ctrl := NewController(inverter, vehicleCtrl, evse, baseConfig(0), metrics)

	// when
	ctrl.Tick()

	// then
	if len(vehicleCtrl.setAmpsCalls) != 0 {
		t.Errorf("Expected no SetChargingAmps calls, got %v", vehicleCtrl.setAmpsCalls)
	}
	if len(metrics.LatestFrame.VehicleFrames) != 0 {
		t.Errorf("Expected no vehicle frames in metrics, got %v", metrics.LatestFrame.VehicleFrames)
	}

	// and
	if metrics.LatestFrame.PowerProductionWatts != 10000 {
		t.Errorf("Wrong power production in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
	if metrics.LatestFrame.PowerConsumptionWatts != 5000 {
		t.Errorf("Wrong power consumption in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
}

func Test_EVSE_NoVehicleConnected(t *testing.T) {
	// given
	inverter := &mockInverter{state: pv.InverterState{PowerProduction: 10000, PowerConsumption: 5000}}
	vehicleCtrl := &mockVehicleController{chargingStates: map[string]vehicle.ChargingState{}}
	evse := &mockEvse{connected: false}
	metrics := metrics.NewMetricsRegistry()
	ctrl := NewController(inverter, vehicleCtrl, evse, baseConfig(0), metrics)

	// when
	ctrl.Tick()

	// then
	if len(vehicleCtrl.setAmpsCalls) != 0 {
		t.Errorf("Expected no SetChargingAmps calls, got %v", vehicleCtrl.setAmpsCalls)
	}
	if len(metrics.LatestFrame.VehicleFrames) != 0 {
		t.Errorf("Expected no vehicle frames in metrics, got %v", metrics.LatestFrame.VehicleFrames)
	}

	// and
	if metrics.LatestFrame.PowerProductionWatts != 10000 {
		t.Errorf("Wrong power production in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
	if metrics.LatestFrame.PowerConsumptionWatts != 5000 {
		t.Errorf("Wrong power consumption in metrics %v", metrics.LatestFrame.PowerProductionWatts)
	}
}

func Test_VehicleOrderPersistsBetweenTicks(t *testing.T) {
	// given
	inverter := &mockInverter{state: pv.InverterState{PowerProduction: 10000, PowerConsumption: 5000}}
	vehicleCtrl := &mockVehicleController{
		chargingStates: map[string]vehicle.ChargingState{
			"VIN3": {Amps: 8, Power: 5},
		},
	}
	evse := &mockEvse{connected: true}
	metrics := metrics.NewMetricsRegistry()
	ctrl := NewController(inverter, vehicleCtrl, evse, threeCarConfig(), metrics)

	assertVehicleOrder(t, ctrl.Vehicles, []string{"VIN1", "VIN2", "VIN3"})

	// when
	ctrl.Tick()

	// then
	assertVehicleOrder(t, ctrl.Vehicles, []string{"VIN3", "VIN1", "VIN2"})

	// when
	vehicleCtrl.getChargingStateCalls = nil
	vehicleCtrl.chargingStates = map[string]vehicle.ChargingState{
		"VIN2": {Amps: 8, Power: 5},
	}
	ctrl.Tick()

	// then
	assertVehicleOrder(t, ctrl.Vehicles, []string{"VIN2", "VIN3", "VIN1"})

	// when
	vehicleCtrl.getChargingStateCalls = nil
	ctrl.Tick()

	// then
	assertVehicleOrder(t, ctrl.Vehicles, []string{"VIN2", "VIN3", "VIN1"})
	if got := vehicleCtrl.getChargingStateCalls; len(got) != 1 || got[0] != "VIN2" {
		t.Errorf("Expected only VIN2 to be checked after it moved to the front, got %v", got)
	}
}

func Test_ManualModes(t *testing.T) {
	for _, tt := range []struct {
		name      string
		mode      Mode
		connected bool
		evseErr   error
		wantMode  Mode
	}{
		{"manual connected", ModeManual, true, nil, ModeManual},
		{"manual disconnected", ModeManual, false, nil, ModeManual},
		{"until unplugged connected", ModeManualUntilDisconnect, true, nil, ModeManualUntilDisconnect},
		{"until unplugged disconnected", ModeManualUntilDisconnect, false, nil, ModePVAutomatic},
		{"until unplugged unknown connection", ModeManualUntilDisconnect, false, errors.New("EVSE unavailable"), ModeManualUntilDisconnect},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// given
			inverter := &mockInverter{state: pv.InverterState{PowerProduction: 10000, PowerConsumption: 5000}}
			vehicleCtrl := &mockVehicleController{chargingStates: map[string]vehicle.ChargingState{"VIN2": {Amps: 8, Power: 5}}}
			evse := &mockEvse{connected: tt.connected, err: tt.evseErr}
			registry := metrics.NewMetricsRegistry()
			ctrl := NewController(inverter, vehicleCtrl, evse, baseConfig(0), registry)
			ctrl.ChangeMode(tt.mode)

			// when
			ctrl.Tick()

			// then
			if ctrl.Mode() != tt.wantMode {
				t.Errorf("Expected mode %v, got %v", tt.wantMode, ctrl.Mode())
			}
			if len(vehicleCtrl.getChargingStateCalls) != 0 || len(vehicleCtrl.setAmpsCalls) != 0 {
				t.Fatal("Manual modes and disconnected vehicles must not query or control vehicles")
			}
			if registry.LatestFrame.PowerProductionWatts != 10000 || registry.LatestFrame.PowerConsumptionWatts != 5000 {
				t.Fatal("Manual modes must continue updating inverter metrics")
			}
		})
	}
}

func Test_ManualUntilDisconnectReturnsToAutomatic(t *testing.T) {
	// given
	// An inverter error must not prevent returning to automatic on disconnect.
	inverter := &mockInverter{err: errors.New("inverter unavailable")}
	vehicleCtrl := &mockVehicleController{chargingStates: map[string]vehicle.ChargingState{"VIN2": {Amps: 8, Power: 5}}}
	evse := &mockEvse{connected: true}
	ctrl := NewController(inverter, vehicleCtrl, evse, baseConfig(0), metrics.NewMetricsRegistry())
	ctrl.ChangeMode(ModeManualUntilDisconnect)
	ctrl.Tick()

	// when
	evse.connected = false
	ctrl.Tick()

	// then
	if ctrl.Mode() != ModePVAutomatic {
		t.Fatal("Disconnect must return to automatic mode")
	}

	// when
	inverter.err = nil
	inverter.state = pv.InverterState{PowerProduction: 10000, PowerConsumption: 5000}
	evse.connected = true
	ctrl.Tick()

	// then
	if ctrl.Mode() != ModePVAutomatic || vehicleCtrl.setAmpsCalls["VIN2"] != 15 {
		t.Fatal("Automatic charging control must resume for the next connected vehicle")
	}
}

// Mocks
type mockEvse struct {
	connected bool
	err       error
}

func (m *mockEvse) IsVehicleConnected() (bool, error) {
	return m.connected, m.err
}

type mockInverter struct {
	state pv.InverterState
	err   error
}

func (m *mockInverter) ReadInverterState() (pv.InverterState, error) {
	return m.state, m.err
}

type mockVehicleController struct {
	chargingStates        map[string]vehicle.ChargingState
	getChargingStateCalls []string
	setAmpsCalls          map[string]int
	err                   error
}

func (m *mockVehicleController) GetChargingState(v *vehicle.Vehicle) (*vehicle.ChargingState, error) {
	m.getChargingStateCalls = append(m.getChargingStateCalls, v.Vin)
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

func assertVehicleOrder(t *testing.T, vehicles []vehicle.Vehicle, expectedVins []string) {
	t.Helper()

	if len(vehicles) != len(expectedVins) {
		t.Fatalf("Expected %d vehicles, got %d", len(expectedVins), len(vehicles))
	}
	for i, vin := range expectedVins {
		if vehicles[i].Vin != vin {
			t.Errorf("Expected vehicle at index %d to be %s, got %s", i, vin, vehicles[i].Vin)
		}
	}
}

func threeCarConfig() *configuration.Config {
	config := baseConfig(0)
	config.Vehicles.Cars = []configuration.Vehicle{
		{Name: "Car1", Vin: "VIN1"},
		{Name: "Car2", Vin: "VIN2"},
		{Name: "Car3", Vin: "VIN3"},
	}
	return config
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
