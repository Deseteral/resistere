package inverter

type Inverter interface {
	ReadEnergySurplus() (energySurplus float64, error error)
}
