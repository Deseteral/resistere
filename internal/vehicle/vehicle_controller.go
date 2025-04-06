package vehicle

type Controller interface {
	SetChargingAmps(vehicle *Vehicle, chargingAmps int) error
}
