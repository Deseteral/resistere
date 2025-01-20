import os.path
import sys
import csv
from datetime import datetime
from pysolarmanv5 import PySolarmanV5

from src.resistere_config import ResistereConfig

power_registers = [
    (0x05C4, "PV Power"),

    (0x0485, "ActivePower_Output_Total"),
    (0x0486, "ReactivePower_Output_Total"),
    (0x0487, "ApparentPower_Output_Total"),

    (0x0488, "ActivePower_PCC_Total"),
    (0x0489, "ReactivePower_PCC_Total"),
    (0x048A, "ApparentPower_PCC_Total"),

    (0x048F, "ActivePower_Output_R"),
    (0x0490, "ReactivePower_Output_R"),

    (0x0493, "ActivePower_PCC_R"),
    (0x0494, "ReactivePower_PCC_R"),

    (0x049A, "ActivePower_Output_S"),
    (0x049B, "ReactivePower_Output_S"),

    (0x049E, "ActivePower_PCC_S"),
    (0x049F, "ReactivePower_PCC_S"),

    (0x04A5, "ActivePower_Output_T"),
    (0x04A6, "ReactivePower_Output_T"),

    (0x04A9, "ActivePower_PCC_T"),
    (0x04AA, "ReactivePower_PCC_T"),

    (0x04AE, "ActivePower_PV_Ext"),
    (0x04AF, "ActivePower_Load_Sys"),

    (0x0504, "ActivePower_Load_Total_EPS"),
    (0x0505, "ReactivePower_Load_Total_EPS"),
    (0x0506, "ApparentPower_Load_Total_EPS"),

    (0x050C, "ActivePower_Load_R_EPS"),
    (0x050D, "ReactivePower_Load_R_EPS"),
    (0x050E, "ApparentPower_Load_R_EPS"),

    (0x0514, "ActivePower_Load_S_EPS"),
    (0x0515, "ReactivePower_Load_S_EPS"),
    (0x0516, "ApparentPower_Load_S_EPS"),

    (0x051C, "ActivePower_Load_T_EPS"),
    (0x051D, "ReactivePower_Load_T_EPS"),
    (0x051E, "ApparentPower_Load_T_EPS"),
]


def persist_values_to_file(config: ResistereConfig):
    inverter = PySolarmanV5(config.inverter.ip, config.inverter.serial, port=config.inverter.port)
    values = read_power_registers(inverter)

    values = [datetime.now().isoformat(), *values]

    filepath = "data.csv"
    file_is_empty = not os.path.exists(filepath)

    with open(filepath, "a") as f:
        w = csv.writer(f)

        if file_is_empty:
            headers = ["date", *map(lambda e: e[1], power_registers)]
            w.writerow(headers)

        w.writerow(values)

    inverter.disconnect()


def read_power_registers(inverter: PySolarmanV5) -> list[int]:
    data = []
    for address, name in power_registers:
        value = read_single_address(address, inverter)
        data.append(value)

        print(f"register: 0x{address:04x}  value: {value}  ({name})")

    return data


def read_single_address(address: int, inverter: PySolarmanV5) -> int:
    try:
        return inverter.read_holding_registers(address, 1)[0]
    except Exception as err:
        print(f"Error while reading value from register 0x{address:04x}.")
        print(err)

        inverter.disconnect()
        sys.exit()
