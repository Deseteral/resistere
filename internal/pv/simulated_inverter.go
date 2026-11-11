package pv

import "sync"

type SimulatedInverter struct {
	mu          sync.RWMutex
	production  float64
	consumption float64
}

func (i *SimulatedInverter) ReadInverterState() (InverterState, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return InverterState{PowerProduction: i.production, PowerConsumption: i.consumption}, nil
}

func (i *SimulatedInverter) SetState(state InverterState) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.production = state.PowerProduction
	i.consumption = state.PowerConsumption
}

func NewSimulatedInverter() *SimulatedInverter {
	return &SimulatedInverter{
		production:  11.0 * 1000.0,
		consumption: 4.0 * 1000.0,
	}
}
