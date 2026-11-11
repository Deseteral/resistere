package webapp

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/evse"
	"github.com/deseteral/resistere/internal/metrics"
	"github.com/deseteral/resistere/internal/pv"
	"github.com/deseteral/resistere/internal/vehicle"
)

func simulatorFixture(t *testing.T, enabled bool) (http.Handler, *Simulator, *controller.Controller, *metrics.Registry) {
	t.Helper()
	s := &Simulator{Inverter: pv.NewSimulatedInverter(), EVSE: evse.NewSimulatedEvse(), Vehicles: vehicle.NewSimulatedVehicleController()}
	config := &configuration.Config{
		SimulatorMode: enabled,
		Controller:    configuration.Controller{GridVoltage: 230, SafetyMarginWatts: 1000},
		Vehicles: configuration.Vehicles{Cars: []configuration.Vehicle{
			{Name: "First car", Vin: "VIN1"}, {Name: "Second car", Vin: "VIN2"},
		}},
	}
	m := metrics.NewMetricsRegistry()
	c := controller.NewController(s.Inverter, s.Vehicles, s.EVSE, config, m)
	router, err := newRouter(config, &c, m, s)
	if err != nil {
		t.Fatal(err)
	}
	return router, s, &c, m
}

func simulatorRequest(router http.Handler, method string, values url.Values) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "/simulator", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}

func simulatorValues() url.Values {
	return url.Values{
		"production": {"7200.5"}, "consumption": {"3100.5"}, "connected": {"on"},
		"vehicle_0_amps": {"9"}, "vehicle_0_power": {"6"},
		"vehicle_1_reachable": {"on"}, "vehicle_1_amps": {"7"}, "vehicle_1_power": {"4"},
	}
}

func TestSimulatorOnlyAvailableInSimulatorMode(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		router, s, _, _ := simulatorFixture(t, enabled)
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			w := simulatorRequest(router, method, simulatorValues())
			want := http.StatusNotFound
			if enabled {
				want = http.StatusOK
			}
			if w.Code != want {
				t.Fatalf("simulator mode %v, %s: got %d, want %d", enabled, method, w.Code, want)
			}
		}
		if !enabled {
			state, _ := s.Inverter.ReadInverterState()
			if state.PowerProduction != 11000 {
				t.Fatal("disabled simulator accepted changes")
			}
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
		if w.Code != http.StatusOK || strings.Contains(w.Body.String(), "href=\"/simulator\"") != enabled {
			t.Fatalf("simulator mode %v: incorrect home page link", enabled)
		}
	}
}

func TestSimulatorParametersDriveController(t *testing.T) {
	for _, scenario := range []struct {
		name      string
		connected bool
		amps      string
		wantAmps  int
	}{
		{"reachable second car gains four amps", true, "7", 11},
		{"disconnected EVSE does not adjust current", false, "7", 7},
		{"idle car does not start charging", true, "0", 0},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			router, s, c, m := simulatorFixture(t, true)
			values := simulatorValues()
			values.Set("vehicle_1_amps", scenario.amps)
			if !scenario.connected {
				values.Del("connected")
			}
			w := simulatorRequest(router, http.MethodPost, values)
			if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "parameters applied") || strings.Contains(w.Body.String(), "<html") {
				t.Fatalf("expected htmx success fragment, got %d: %s", w.Code, w.Body.String())
			}
			if _, err := s.Vehicles.GetChargingState(&vehicle.Vehicle{Vin: "VIN1"}); err == nil {
				t.Fatal("first car should be unreachable")
			}
			c.Tick()
			if got := s.Vehicles.SimulationState("VIN2"); got.Amps != scenario.wantAmps || got.Power != 4 {
				t.Fatalf("second car state = %+v, want %d A and 4 kW", got, scenario.wantAmps)
			}
			if got := s.Vehicles.SimulationState("VIN1"); got.Amps != 9 || got.Power != 6 {
				t.Fatalf("first car unexpectedly changed: %+v", got)
			}
			if m.LatestFrame.PowerProductionWatts != 7200.5 || m.LatestFrame.PowerConsumptionWatts != 3100.5 {
				t.Fatalf("inverter parameters did not reach metrics: %+v", m.LatestFrame)
			}
			if scenario.connected && scenario.wantAmps > 0 && (len(m.LatestFrame.VehicleFrames) != 1 || m.LatestFrame.VehicleFrames[0].ChargingPowerWatts != 4000) {
				t.Fatalf("reported charging power did not reach metrics: %+v", m.LatestFrame)
			}
			w = simulatorRequest(router, http.MethodGet, nil)
			if !strings.Contains(w.Body.String(), "value=\""+strconv.Itoa(scenario.wantAmps)+"\"") {
				t.Fatal("refresh did not include controller's current")
			}
		})
	}
}

func TestSimulatorRejectsInvalidParametersWithoutPartialUpdates(t *testing.T) {
	for field, invalid := range map[string][]string{
		"production": {"NaN", "+Inf", "-1", ""}, "consumption": {"NaN", "-2"},
		"vehicle_1_amps": {"17", "-1", "1.5"}, "vehicle_1_power": {"12", "-1", "oops"},
	} {
		for _, value := range invalid {
			t.Run(field+"="+value, func(t *testing.T) {
				router, s, _, _ := simulatorFixture(t, true)
				values := simulatorValues()
				values.Set(field, value)
				values.Del("connected")
				w := simulatorRequest(router, http.MethodPost, values)
				if !strings.Contains(w.Body.String(), "No parameters changed.") {
					t.Fatal("missing validation message")
				}
				state, _ := s.Inverter.ReadInverterState()
				connected, _ := s.EVSE.IsVehicleConnected()
				if state.PowerProduction != 11000 || state.PowerConsumption != 4000 || !connected || s.Vehicles.SimulationState("VIN1").Amps != 8 {
					t.Fatal("invalid form partially updated simulator")
				}
			})
		}
	}
}
