package vehicle

import (
	"errors"
	"sync"
)

type SimulatedVehicleController struct {
	mu     sync.RWMutex
	states map[string]SimulatedVehicleState
}

type SimulatedVehicleState struct {
	Reachable bool
	ChargingState
}

func (c *SimulatedVehicleController) GetChargingState(vehicle *Vehicle) (state *ChargingState, error error) {
	s := c.SimulationState(vehicle.Vin)
	if !s.Reachable {
		return nil, errors.New("simulated vehicle is out of range")
	}
	return &s.ChargingState, nil
}

func (c *SimulatedVehicleController) SetChargingAmps(vehicle *Vehicle, chargingAmps int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := c.simulationState(vehicle.Vin)
	s.Amps = chargingAmps
	c.states[vehicle.Vin] = s
	return nil
}

func (c *SimulatedVehicleController) SimulationState(vin string) SimulatedVehicleState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.simulationState(vin)
}

func (c *SimulatedVehicleController) simulationState(vin string) SimulatedVehicleState {
	if state, ok := c.states[vin]; ok {
		return state
	}
	return SimulatedVehicleState{Reachable: true, ChargingState: ChargingState{Amps: 8, Power: 5}}
}

func (c *SimulatedVehicleController) SetSimulationState(vin string, state SimulatedVehicleState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states[vin] = state
}

func NewSimulatedVehicleController() *SimulatedVehicleController {
	return &SimulatedVehicleController{
		states: make(map[string]SimulatedVehicleState),
	}
}
