package vehicle

type SimulatedVehicleController struct{}

func (c *SimulatedVehicleController) GetChargingAmps(vehicle *Vehicle) (amps int, error error) {
	return 16, nil
}

func (c *SimulatedVehicleController) SetChargingAmps(vehicle *Vehicle, chargingAmps int) error {
	return nil
}

func NewSimulatedVehicleController() *SimulatedVehicleController {
	return &SimulatedVehicleController{}
}
