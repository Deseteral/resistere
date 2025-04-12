package vehicle

type SimulatedVehicleController struct{}

func (c *SimulatedVehicleController) SetChargingAmps(vehicle *Vehicle, chargingAmps int) error {
	return nil
}

func NewSimulatedVehicleController() *SimulatedVehicleController {
	return &SimulatedVehicleController{}
}
