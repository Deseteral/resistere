package vehicle

type Controller interface {
	GetChargingState(vehicle *Vehicle) (state *ChargingState, error error)
	SetChargingAmps(vehicle *Vehicle, chargingAmps int) error
}

type ChargingState struct {
	Amps  int
	Power int
}
