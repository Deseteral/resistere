package pv

type Inverter interface {
	ReadInverterState() (InverterState, error)
}

type InverterState struct {
	PowerProduction  float64
	PowerConsumption float64
}
