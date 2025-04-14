package pv

import (
	"bytes"
	"embed"
	"github.com/deseteral/resistere/internal/configuration"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// TODO: Ideally this would not require using Python.
//       This may or may not require rewriting significant portions of pysolarmanv5 library to Go lang.

//go:embed solarman_interface/build/solarman_interface.pyz
var solarmanInterfaceBinary embed.FS

type SolarmanInverter struct {
	ip     string
	serial string
	port   string
}

func (i *SolarmanInverter) ReadEnergySurplus() (InverterState, error) {
	// Extract Python binary to tmp location for running.
	binaryFilePath, err := i.preparePythonBinary()
	if err != nil {
		return InverterState{}, err
	}

	// Remove Python binary from tmp location after we're done reading.
	defer func(binaryFilePath string) {
		// We can ignore the error here as it's not crucial to clean up.
		// This might eventually become a problem, so logging that this happened is worth doing.
		_ = i.cleanupSolarmanInterface(binaryFilePath)
	}(binaryFilePath)

	// Process and return output from Python binary.
	state, err := i.execPythonBinary(binaryFilePath)
	if err != nil {
		return InverterState{}, err
	}

	log.Printf("Read from Solarman inverter: production %fkW, consuption %fkW.\n", state.PowerProduction, state.PowerConsumption)

	return state, nil
}

func (i *SolarmanInverter) preparePythonBinary() (binaryPath string, error error) {
	tmpFile, err := os.CreateTemp("", "solarman_interface_*.pyz")
	if err != nil {
		log.Printf("Failed to create temporary file for solarman_interface: %v", err)
		return "", err
	}

	data, err := solarmanInterfaceBinary.ReadFile("solarman_interface/build/solarman_interface.pyz")
	if err != nil {
		log.Printf("Failed to read solarman_interface binary: %v", err)
		return "", err
	}

	_, err = tmpFile.Write(data)
	if err != nil {
		log.Printf("Failed to write solarman_interface to tmp file: %v", err)
		return "", err
	}

	err = tmpFile.Close()
	if err != nil {
		log.Printf("Failed to close solarman_interface tmp file: %v", err)
		return "", err
	}

	err = os.Chmod(tmpFile.Name(), 0777)
	if err != nil {
		log.Printf("Failed to set permiossions for solarman_interface binary: %v", err)
		return "", err
	}

	return tmpFile.Name(), nil
}

func (i *SolarmanInverter) cleanupSolarmanInterface(binaryFilePath string) error {
	err := os.Remove(binaryFilePath)
	if err != nil {
		log.Printf("Failed to remove solarman_interface tmp file: %v", err)
		return err
	}
	return nil
}

func (i *SolarmanInverter) execPythonBinary(binaryFilePath string) (InverterState, error) {
	cmd := exec.Command(binaryFilePath, i.ip, i.serial, i.port)

	var buffer bytes.Buffer
	cmd.Stdout = &buffer

	err := cmd.Run()
	if err != nil {
		log.Printf("solarman_interface failed to produce expeted output: %v", err)
		return InverterState{}, err
	}

	output := buffer.String()
	values := strings.Split(output, " ")

	powerProduction, err := strconv.ParseFloat(strings.TrimSpace(values[0]), 64)
	if err != nil {
		log.Printf("Could not parse solarman_interface output: %v", err)
		return InverterState{}, err
	}

	powerConsumption, err := strconv.ParseFloat(strings.TrimSpace(values[1]), 64)
	if err != nil {
		log.Printf("Could not parse solarman_interface output: %v", err)
		return InverterState{}, err
	}

	state := InverterState{PowerProduction: powerProduction, PowerConsumption: powerConsumption}
	return state, nil
}

func NewSolarmanInverter(config *configuration.SolarmanInverter) *SolarmanInverter {
	return &SolarmanInverter{
		ip:     config.Ip,
		serial: config.Serial,
		port:   config.Port,
	}
}
