package vehicle

type SimulatedVehicleController struct {
	chargingAmps int
	power        int
}

func (c *SimulatedVehicleController) GetChargingState(vehicle *Vehicle) (state *ChargingState, error error) {
	return &ChargingState{Amps: c.chargingAmps, Power: 5}, nil
}

func (c *SimulatedVehicleController) SetChargingAmps(vehicle *Vehicle, chargingAmps int) error {
	c.chargingAmps = chargingAmps
	return nil
}

func NewSimulatedVehicleController() *SimulatedVehicleController {
	return &SimulatedVehicleController{
		chargingAmps: 8,
		power:        5,
	}
}
