package evse

type Evse interface {
	// IsVehicleConnected TODO: See if we need this after PoC is up and running. If it's not needed - remove.
	IsVehicleConnected() (isVehicleConnected bool, error error)
}
