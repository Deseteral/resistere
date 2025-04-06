import argparse
import sys

from pysolarmanv5 import PySolarmanV5


def main():
    parser = argparse.ArgumentParser(description="Read current energy surplus from Solarman Inverter.")
    parser.add_argument("ip", type=str, help="IP address of the inverter.")
    parser.add_argument("serial", type=int, help="Serial number of the inverter.")
    parser.add_argument("port", type=int, help="Port number of the inverter.")
    args = parser.parse_args()

    surplus = read_current_energy_surplus(args.ip, args.serial, args.port)

    if surplus is not None:
        print(surplus)
        sys.exit(0)
    else:
        sys.exit(1)


def read_current_energy_surplus(ip: str, serial: int, port: int) -> float | None:
    inverter = _connect_to_inverter(ip, serial, port)
    if inverter is None:
        return None

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


def _connect_to_inverter(ip: str, serial: int, port: int) -> PySolarmanV5 | None:
    try:
        inverter = PySolarmanV5(ip, serial, port=port)
        return inverter
    except Exception as err:
        print("Error while connecting to inverter.")
        print(err)
        return None


if __name__ == "__main__":
    main()
