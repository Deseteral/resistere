package evse

import "sync"

type SimulatedEvse struct {
	mu        sync.RWMutex
	connected bool
}

func (e *SimulatedEvse) IsVehicleConnected() (isVehicleConnected bool, error error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.connected, nil
}

func (e *SimulatedEvse) SetConnected(connected bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.connected = connected
}

func NewSimulatedEvse() *SimulatedEvse {
	return &SimulatedEvse{connected: true}
}
