# resistere

This project provides a Raspberry Pi-based controller designed to dynamically adjust the charging current for Tesla vehicles.

It aims to maximize the use of solar energy by ensuring the car charges only from PV production surplus.

## Controller logic

Core of processing logic is contained in `internal/controller.go` module. Following diagram presents simplified cycle flow.

```
┌───────────────┐             ┌───────────────┐        ┌ ─ ─ ─ ─ ─ ─ ─ ┐      ┌───────────────┐
│   resistere   │░            │  PV inverter  │░           Vehicle 1          │   Vehicle 2   │░
└───────────────┘░            └───────────────┘░       └ ─ ─ ─ ─ ─ ─ ─ ┘      └───────────────┘░
 ░░░░░░░│░░░░░░░░░             ░░░░░░░│░░░░░░░░░               │               ░░░░░░░│░░░░░░░░░

        │ Get current power production│                        │                      │
                and consumption
        │───────────────────────────▶┌┴┐                       │                      │
                                     │ │
        │                            │ │                       │                      │
         ◀───────────────────────────└─┘
        │                             │                        │                      │
       ┌─┐───┐
       │ │   │Calculate energy        │                        │                      │
       │ │   │surplus
       └┬┘◀──┘                        │                        │                      │

        │                   Is vehicle in range?               │                      │
         ────────────────────────────────────────────────────▶┌─┐
        │                             │                       │ │                     │
                                                              │ │
        │                    Connection timeout               │ │                     │
         ◀────────────────────────────────────────────────────└─┘
        │                             │                        │                      │
                                         Is vehicle in range?
        │─────────────────────────────┼────────────────────────┼────────────────────▶┌┴┐
                                                                                     │ │
        │                             │                        │                     │ │
                                        In range and charging                        │ │
        │◀────────────────────────────┼────────────────────────┼─────────────────────└┬┘

        │                             │                        │                      │
       ┌─┐───┐
       │ │   │Calculate 𝚫 amps        │                        │                      │
       │ │   │
       └┬┘◀──┘                        │                        │                      │

        │                             │  Set charging amps     │                      │
         ────────────────────────────────────────────────────────────────────────────▶
        │                             │                        │                      │

┌───────────────┐             ┌───────────────┐        ┌ ─ ─ ─ ─ ─ ─ ─ ┐      ┌───────────────┐
│   resistere   │░            │  PV inverter  │░           Vehicle 1          │   Vehicle 2   │░
└───────────────┘░            └───────────────┘░       └ ─ ─ ─ ─ ─ ─ ─ ┘      └───────────────┘░
 ░░░░░░░░░░░░░░░░░             ░░░░░░░░░░░░░░░░░                               ░░░░░░░░░░░░░░░░░
```

To calculate delta amps, controller has to calculate current power surplus (difference between power production and consumption) and read current charging amps from vehicle. Then calculates it using following formula.

```
𝚫A = W / (V * 3)
 ▲   ▲    ▲   ▲
 │   │    │   └──EVSE is using three-phases to charge.
 │   │    └──the electric potential of energy grid.
 │   └──the energy surplus (the difference between power generated by PV and total power used).
 └──the amount by which we can increase or decrease charging current.
```

## Hardware

I've made this for myself and while it's fairly configurable, it does require specific hardware combination to work:

- Any Tesla vehicle that can be controlled via Bluetooth.
- Sofar HYD 5-20KTL-3PH inverter with Wi-Fi data logger (but it should work with any inverter that outputs the Modbus data in similar way).
- Linux-based device with Bluetooth Low Energy connectivity (like Raspberry Pi Zero 2 W).

Adding support for other inverters should be easy, provided you can connect to them and read current energy production and consumption.

## Configuration

The application requires `config.toml` file to run. Here's example configuration:

```toml
[web]
port = 5467 # Port on which the web UI will be running.

[controller]
cycle_interval_seconds = 10 # Time between charging current modulation (in seconds).
safety_margin_watts = 1000  # Surplus power safety margin (in watts).
grid_voltage = 230          # Grid voltage in your country/area.

[solarman_inverter]
ip = "192.168.1.99"   # IP address of inverter's data logger in your local network.
serial = "1234567890" # Serial number of inverter's data logger.
port = "8899"         # Port of inverter's data logger.

[tesla_control]
key_file = "./private_key.pem" # Private key paired with your Teslas.

[vehicles]
# List of Tesla vehicles to connect to (name is only needed for easier identification in logs).
cars = [
    { name = "My Model 3", vin = "QWERTYUIOP1234567" },
    { name = "My Model Y", vin = "ABCDEFGHIJ1234567" },
]
```

## 📝 License

This project is licensed under the [MIT license](LICENSE).
