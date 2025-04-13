package vehicle

type Controller interface {
	IsCharging(vehicle *Vehicle) (isCharging bool, error error)
	SetChargingAmps(vehicle *Vehicle, chargingAmps int) error
}
