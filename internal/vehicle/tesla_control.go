package vehicle

import (
	"github.com/deseteral/resistere/internal/configuration"
	"os/exec"
	"strconv"
)

type TeslaControlController struct {
	keyFilePath string
}

func (c TeslaControlController) SetChargingAmps(vehicle *Vehicle, amps int) error {
	cmd := exec.Command(
		"tesla-control",
		"-vin", vehicle.Vin,
		"-key-file", c.keyFilePath,
		"-ble",
		"charging-set-amps", strconv.Itoa(amps),
	)
	return cmd.Run()
}

func NewTeslaControlController(config *configuration.TeslaControl) TeslaControlController {
	return TeslaControlController{
		keyFilePath: config.KeyFile,
	}
}
