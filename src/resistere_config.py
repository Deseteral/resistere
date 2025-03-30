import tomllib
from dataclasses import dataclass


@dataclass
class InverterConfig:
    enabled: bool
    ip: str
    serial: int
    port: int


@dataclass
class ControllerConfig:
    cycle_interval_seconds: int


@dataclass
class ResistereConfig:
    inverter: InverterConfig
    controller: ControllerConfig


def read_config() -> ResistereConfig:
    with open("./resistere_config.toml", "rb") as f:
        data = tomllib.load(f)

        config = ResistereConfig(
            InverterConfig(data["inverter"]["enabled"], data["inverter"]["ip"], data["inverter"]["serial"], data["inverter"]["port"]),
            ControllerConfig(data["controller"]["cycle_interval_seconds"])
        )

        return config
