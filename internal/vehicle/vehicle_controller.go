package vehicle

type Controller interface {
	GetChargingAmps(vehicle *Vehicle) (amps int, error error)
	SetChargingAmps(vehicle *Vehicle, chargingAmps int) error
}
