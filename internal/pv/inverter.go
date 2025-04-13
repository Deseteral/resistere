package pv

type Inverter interface {
	ReadEnergySurplus() (InverterState, error)
}

type InverterState struct {
	PowerProduction  float64
	PowerConsumption float64
}
