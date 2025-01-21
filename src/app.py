from time import sleep

from src import pv
from src.resistere_config import ResistereConfig


def run(config: ResistereConfig):
    print("Starting Resistere controller.")

    print("Configuration:")
    print(config)

    while True:
        _cycle(config)
        sleep(config.controller.cycle_interval_seconds)


def _cycle(config: ResistereConfig):
    surplus_power = pv.read_current_energy_surplus(config)
    if surplus_power is None:
        return

    print(f"Current energy surplus: {surplus_power} kW")
