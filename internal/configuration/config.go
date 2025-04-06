package configuration

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	Web              Web
	Controller       Controller
	SolarmanInverter SolarmanInverter `toml:"solarman_inverter"`
	TeslaControl     TeslaControl     `toml:"tesla_control"`
	SimulatorMode    bool             `toml:"simulator_mode"`
}

type Web struct {
	Port int
}

type Controller struct {
	CycleIntervalSeconds int `toml:"cycle_interval_seconds"`
}

type SolarmanInverter struct {
	Ip     string
	Serial string
	Port   string
}

type TeslaControl struct {
	KeyFile string `toml:"key_file"`
}

func ReadConfig() (*Config, error) {
	file, err := os.ReadFile("resistere_config.toml")
	if err != nil {
		log.Printf("Could not read config file: %v", err)
		return nil, err
	}

	tomlData := string(file)

	var conf Config
	_, err = toml.Decode(tomlData, &conf)
	if err != nil {
		log.Printf("Could not parse config file: %v", err)
		return nil, err
	}

	log.Println("Configuration successfully read.")
	return &conf, err
}
