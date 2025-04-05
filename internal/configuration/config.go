package configuration

import (
    "github.com/BurntSushi/toml"
    "log"
    "os"
)

type Config struct {
    Web Web
}

type Web struct {
    Port int
}

func ReadConfig() (*Config, error) {
    file, err := os.ReadFile("resistere_config.toml")
    if err != nil {
        return nil, err
    }

    var tomlData = string(file)

    var conf Config
    _, err = toml.Decode(tomlData, &conf)
    if err != nil {
        return nil, err
    }

    log.Println("Configuration successfully read.")
    return &conf, err
}
