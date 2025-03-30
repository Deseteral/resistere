from src import resistere_config
from src import app
from src.app import App
from src.solarman_inverter import SolarmanInverter


def main():
    config = resistere_config.read_config()
    inverter = SolarmanInverter(config)
    app = App(config, inverter)

    app.run()


if __name__ == "__main__":
    main()
