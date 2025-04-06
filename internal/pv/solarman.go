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

//go:embed solarman_interface/build/solarman_interface.pyz
var solarmanInterfaceBinary embed.FS

type SolarmanInverter struct {
	Ip     string
	Serial string
	Port   string
}

func (i SolarmanInverter) ReadEnergySurplus() (energySurplus float64, error error) {
	// Extract Python binary to tmp location for running.
	binaryFilePath, err := i.preparePythonBinary()
	if err != nil {
		return -1, err
	}

	// Remove Python binary from tmp location after we're done reading.
	defer func(binaryFilePath string) {
		// We can ignore the error here as it's not crucial to clean up.
		// This might eventually become a problem, so logging that this happened is worth doing.
		_ = i.cleanupSolarmanInterface(binaryFilePath)
	}(binaryFilePath)

	// Process and return output from Python binary.
	energySurplus, err = i.execPythonBinary(binaryFilePath)
	if err != nil {
		return -1, err
	}

	return energySurplus, nil
}

func (i SolarmanInverter) preparePythonBinary() (binaryPath string, error error) {
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

	return tmpFile.Name(), nil
}

func (i SolarmanInverter) cleanupSolarmanInterface(binaryFilePath string) error {
	err := os.Remove(binaryFilePath)
	if err != nil {
		log.Printf("Failed to remove solarman_interface tmp file: %v", err)
		return err
	}
	return nil
}

func (i SolarmanInverter) execPythonBinary(binaryFilePath string) (energySurplus float64, error error) {
	cmd := exec.Command(binaryFilePath, i.Ip, i.Serial, i.Port)

	var buffer bytes.Buffer
	cmd.Stdout = &buffer

	err := cmd.Run()
	if err != nil {
		log.Printf("solarman_interface failed to produce expeted output: %v", err)
		return -1, err
	}

	output := buffer.String()

	energySurplus, err = strconv.ParseFloat(strings.TrimSpace(output), 64)
	if err != nil {
		log.Printf("Could not parse solarman_interface output: %v", err)
		return -1, err
	}

	return energySurplus, nil
}

func NewSolarmanInverter(config *configuration.SolarmanInverter) SolarmanInverter {
	return SolarmanInverter{
		Ip:     config.Ip,
		Serial: config.Serial,
		Port:   config.Port,
	}
}
