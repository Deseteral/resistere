""" Scan Modbus registers to find valid registers"""
from pysolarmanv5 import PySolarmanV5
import resistere_config


def main():
    config = resistere_config.read_config()

    modbus = PySolarmanV5(
        config.inverter.ip, config.inverter.serial, port=config.inverter.port
        # , mb_slave_id=1, verbose=False
    )

    # print("Scanning input registers")
    # for x in range(30000, 39999):
    #     try:
    #         val = modbus.read_input_registers(register_addr=x, quantity=1)[0]
    #         print(f"Register: {x:05}\t\tValue: {val:05} ({val:#06x})")
    #     except (V5FrameError, umodbus.exceptions.IllegalDataAddressError):
    #         continue
    # print("Finished scanning input registers")

    print("Scanning holding registers")
    for x in range(1000, 9999):
    # for x in [1476]:
        try:
            val = modbus.read_holding_registers(register_addr=x, quantity=1)[0]
            print(f"Register: {x:05}\t\tValue: {val / 10.0} ({val:#06x})")
        except:
            continue
    print("Finished scanning holding registers")


if __name__ == "__main__":
    main()
