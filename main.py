from src import resistere_config
from src import app
from src.app import App
from src.fake_inverter import FakeInverter
from src.solarman_inverter import SolarmanInverter


def main():
    config = resistere_config.read_config()
    inverter = SolarmanInverter(config) if config.inverter.enabled else FakeInverter()
    app = App(config, inverter)

    app.run()


if __name__ == "__main__":
    main()
