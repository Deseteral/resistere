package configuration

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	SimulatorMode      bool `toml:"simulator_mode"`
	Web                Web
	Controller         Controller
	SolarmanInverter   SolarmanInverter   `toml:"solarman_inverter"`
	TeslaControl       TeslaControl       `toml:"tesla_control"`
	TeslaWallConnector TeslaWallConnector `toml:"tesla_wall_connector"`
	Vehicles           Vehicles
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

type TeslaWallConnector struct {
	Ip string
}

type Vehicles struct {
	Cars []Vehicle
}

type Vehicle struct {
	Name string
	Vin  string
}

func ReadConfig() (*Config, error) {
	file, err := os.ReadFile("config.toml")
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
