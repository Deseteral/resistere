from pysolarmanv5 import PySolarmanV5

from src.resistere_config import ResistereConfig

class SolarmanInverter:
    def __init__(self, config: ResistereConfig):
        self.config = config

    def read_current_energy_surplus(self) -> float | None:
        inverter = self._connect_to_inverter()
        if inverter is None:
            return None

        pv_power = self._read_single_address(0x05C4, inverter) * 0.1
        if pv_power is None:
            return None

        consumption = self._read_single_address(0x04AF, inverter) * 0.01
        if consumption is None:
            return None

        inverter.disconnect()

        surplus = max(pv_power - consumption, 0)
        return surplus


    @staticmethod
    def _read_single_address(address: int, inverter: PySolarmanV5) -> int | None:
        try:
            return inverter.read_holding_registers(address, 1)[0]
        except Exception as err:
            print(f"Error while reading value from register 0x{address:04x}.")
            print(err)
            inverter.disconnect()
            return None

    def _connect_to_inverter(self) -> PySolarmanV5 | None:
        try:
            inverter = PySolarmanV5(self.config.inverter.ip, self.config.inverter.serial, port=self.config.inverter.port)
            return inverter
        except Exception as err:
            print("Error while connecting to inverter.")
            print(err)
            return None
