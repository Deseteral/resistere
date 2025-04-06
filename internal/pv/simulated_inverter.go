package pv

type SimulatedInverter struct{}

func (i SimulatedInverter) ReadEnergySurplus() (energySurplus float64, error error) {
	return 5, nil
}

func NewSimulatedInverter() SimulatedInverter {
	return SimulatedInverter{}
}
