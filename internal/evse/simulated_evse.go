package evse

type SimulatedEvse struct{}

func (e *SimulatedEvse) IsVehicleConnected() (isVehicleConnected bool, error error) {
	return true, nil
}

func NewSimulatedEvse() *SimulatedEvse {
	return &SimulatedEvse{}
}
