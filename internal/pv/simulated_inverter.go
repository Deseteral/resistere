package pv

type SimulatedInverter struct {
	production  float64
	consumption float64
}

func (i *SimulatedInverter) ReadInverterState() (InverterState, error) {
	return InverterState{PowerProduction: i.production, PowerConsumption: i.consumption}, nil
}

func NewSimulatedInverter() *SimulatedInverter {
	return &SimulatedInverter{
		production:  11.0 * 1000.0,
		consumption: 4.0 * 1000.0,
	}
}
