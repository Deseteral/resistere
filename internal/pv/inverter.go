package pv

type Inverter interface {
	// ReadEnergySurplus TODO: Rename.
	ReadEnergySurplus() (InverterState, error)
}

type InverterState struct {
	PowerProduction  float64
	PowerConsumption float64
}
