package vehicle

type SimulatedVehicleController struct {
	chargingAmps int
}

func (c *SimulatedVehicleController) GetChargingAmps(vehicle *Vehicle) (amps int, error error) {
	return c.chargingAmps, nil
}

func (c *SimulatedVehicleController) SetChargingAmps(vehicle *Vehicle, chargingAmps int) error {
	c.chargingAmps = chargingAmps
	return nil
}

func NewSimulatedVehicleController() *SimulatedVehicleController {
	return &SimulatedVehicleController{}
}
