package configuration

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	Web        Web
	Controller Controller
	Inverter   Inverter
}

type Web struct {
	Port int
}

type Controller struct {
	CycleIntervalSeconds int `toml:"cycle_interval_seconds"`
}

type Inverter struct {
	Ip     string
	Serial string
	Port   string
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
