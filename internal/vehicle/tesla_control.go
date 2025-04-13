package vehicle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deseteral/resistere/internal/configuration"
	"os/exec"
	"strconv"
	"strings"
)

// TODO: This implementation is kind of stupid.
//       Instead of invoking tesla-control, we could just integrate with vehicle-command library.
//       This would be much simpler code to handle as it would not require parsing errors and responses from
//       command line interface.

type TeslaControlController struct {
	keyFilePath string
}

func (c *TeslaControlController) GetChargingAmps(vehicle *Vehicle) (amps int, error error) {
	cmd := exec.Command(
		"tesla-control",
		"-vin", vehicle.Vin,
		"-key-file", c.keyFilePath,
		"-ble",
		"-command-timeout", "3s",
		"-connect-timeout", "3s",
		"state",
	)

	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		stderrContent := strings.TrimSuffix(errBuffer.String(), "\n")
		return -1, errors.New(fmt.Sprintf("error while running tesla-control: %s", stderrContent))
	}

	output := outBuffer.String()

	var data map[string]interface{}
	err = json.Unmarshal([]byte(output), &data)
	if err != nil {
		return -1, err
	}

	chargeState, ok := data["chargeState"].(map[string]interface{})
	if !ok {
		return -1, errors.New("error parsing tesla-control state JSON: chargeState is undefined")
	}

	chargingState, ok := chargeState["chargingState"].(map[string]interface{})
	if !ok {
		return -1, errors.New("error parsing tesla-control state JSON: chargeState.chargingState is undefined")
	}

	_, charging := chargingState["Charging"]
	if !charging {
		return -1, nil
	}

	amps = chargeState["chargingAmps"].(int)

	return amps, nil
}

func (c *TeslaControlController) SetChargingAmps(vehicle *Vehicle, amps int) error {
	cmd := exec.Command(
		"tesla-control",
		"-vin", vehicle.Vin,
		"-key-file", c.keyFilePath,
		"-ble",
		"charging-set-amps", strconv.Itoa(amps),
	)
	return cmd.Run()
}

func NewTeslaControlController(config *configuration.TeslaControl) *TeslaControlController {
	return &TeslaControlController{
		keyFilePath: config.KeyFile,
	}
}
