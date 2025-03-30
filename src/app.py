from time import sleep

from src.solarman_inverter import SolarmanInverter
from src.resistere_config import ResistereConfig

class App:
    def __init__(self, config: ResistereConfig, inverter: SolarmanInverter):
        self.config = config
        self.inverter = inverter

    def run(self):
        print("Starting Resistere controller.")

        print("Configuration:")
        print(self.config)

        while True:
            self._cycle()
            sleep(self.config.controller.cycle_interval_seconds)

    def _cycle(self):
        surplus_power = self.inverter.read_current_energy_surplus()
        if surplus_power is None:
            return

        print(f"Current energy surplus: {surplus_power} kW")
