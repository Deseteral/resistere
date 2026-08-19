package evse

type Evse interface {
	IsVehicleConnected() (isVehicleConnected bool, error error)
}
