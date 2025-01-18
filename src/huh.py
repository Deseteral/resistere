from pysolarmanv5 import PySolarmanV5
import time

import resistere_config


def main():
    config = resistere_config.read_config()

    while True:
        # that is minimal code to read the actual PV overall power from Sofar
        # Hybrid usin Solaman whatever stick
        link_to_inverter_using_solarman_stick = PySolarmanV5(
            config.inverter.ip, config.inverter.serial, port=config.inverter.port
        )
        try_asking_inverter = True
        faulty_attempts = 0
        # in case of parallel communication with the inverter, with solarman
        # cloud and your script
        # it sometimes happens you receive frame which was not meant for you,
        # but for the colud
        # then pysolarmanv5 throws an exception - try / except is to hanlde
        # this (and other ) exceptions
        # I suggest to wait second or few and query the inverter again
        while try_asking_inverter:
            try:
                received_list_of_registers = (
                    link_to_inverter_using_solarman_stick.read_holding_registers(1476, 1),
                    link_to_inverter_using_solarman_stick.read_holding_registers(0x0468, 1),
                )
                
                print(received_list_of_registers)

                # actual_PV_overall_power = received_list_of_registers[0] / 10.0
                # print(f"Actual PV overall power: {actual_PV_overall_power} kW")
                #
                # grid_output = received_list_of_registers[1] / 10.0
                # print(f"Output power: {grid_output} kW")

                try_asking_inverter = False
                link_to_inverter_using_solarman_stick.disconnect()
            except Exception as err:
                print(err)
                print("Error in response, trying again in 1 second")
                time.sleep(1)
                faulty_attempts += 1
                if faulty_attempts > 5:
                    print("Too many faulty attempts, exiting")
                    break
        time.sleep(5)


if __name__ == "__main__":
    main()
