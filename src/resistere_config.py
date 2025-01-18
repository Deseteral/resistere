import tomllib
from dataclasses import dataclass

@dataclass
class InverterConfig:
    ip: str
    serial: int
    port: int

@dataclass
class ResistereConfig:
    inverter: InverterConfig

def read_config() -> ResistereConfig:
    with open("./resistere-config.toml", "rb") as f:
        data = tomllib.load(f)

        config = ResistereConfig(
            InverterConfig(data["inverter"]["ip"], data["inverter"]["serial"], data["inverter"]["port"]),
        )

        return config
