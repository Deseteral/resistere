package vehicle

type Controller interface {
	SetChargingAmps(chargingAmps int) error
}
