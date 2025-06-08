package vehicle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/deseteral/resistere/internal/configuration"
)

type TeslaControlController struct {
	keyFilePath string
}

func (c *TeslaControlController) GetChargingState(vehicle *Vehicle) (state *ChargingState, error error) {
	cmd := exec.Command(
		"tesla-control",
		"-vin", vehicle.Vin,
		"-key-file", c.keyFilePath,
		"-ble",
		"-command-timeout", "10s",
		"-connect-timeout", "10s",
		"state", "charge",
	)

	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		stderrContent := strings.TrimSuffix(errBuffer.String(), "\n")
		return nil, errors.New(fmt.Sprintf("error while running tesla-control: %s", stderrContent))
	}

	output := outBuffer.String()

	var data map[string]any
	err = json.Unmarshal([]byte(output), &data)
	if err != nil {
		return nil, err
	}

	chargeState, ok := data["chargeState"].(map[string]any)
	if !ok {
		return nil, errors.New("error parsing tesla-control state JSON: chargeState")
	}

	chargingState, ok := chargeState["chargingState"].(map[string]any)
	if !ok {
		return nil, errors.New("error parsing tesla-control state JSON: chargeState.chargingState")
	}

	_, charging := chargingState["Charging"]
	if !charging {
		return nil, nil
	}

	chargingAmpsRaw, ok := chargeState["chargingAmps"].(float64)
	if !ok {
		return nil, errors.New("error parsing tesla-control state JSON: chargeState.chargingAmps")
	}

	powerRaw, ok := chargeState["chargerPower"].(float64)
	if !ok {
		return nil, errors.New("error parsing tesla-control state JSON: chargeState.chargerPower")
	}

	s := ChargingState{
		Amps:  int(chargingAmpsRaw),
		Power: int(powerRaw),
	}

	return &s, nil
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
