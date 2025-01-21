from pysolarmanv5 import PySolarmanV5

from src.resistere_config import ResistereConfig


def read_current_energy_surplus(config: ResistereConfig) -> float | None:
    inverter = PySolarmanV5(config.inverter.ip, config.inverter.serial, port=config.inverter.port)

    pv_power = _read_single_address(0x05C4, inverter) * 0.1
    if pv_power is None:
        return None

    consumption = _read_single_address(0x04AF, inverter) * 0.01
    if consumption is None:
        return None

    inverter.disconnect()

    surplus = max(pv_power - consumption, 0)
    return surplus


def _read_single_address(address: int, inverter: PySolarmanV5) -> int | None:
    try:
        return inverter.read_holding_registers(address, 1)[0]
    except Exception as err:
        print(f"Error while reading value from register 0x{address:04x}.")
        print(err)
        inverter.disconnect()
        return None
