package pv

type SimulatedInverter struct{}

func (i *SimulatedInverter) ReadEnergySurplus() (InverterState, error) {
	return InverterState{PowerProduction: 11, PowerConsumption: 3}, nil
}

func NewSimulatedInverter() *SimulatedInverter {
	return &SimulatedInverter{}
}
