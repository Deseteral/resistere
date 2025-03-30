class FakeInverter:
    def __init__(self):
        print("Running fake inverter.")

    def read_current_energy_surplus(self) -> float | None:
        return 5.0
