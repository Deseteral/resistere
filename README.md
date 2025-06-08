# resistere

`resistere` is a solution for dynamic power management that adjusts charge rate of Tesla vehicles based on the PV production surplus.

It aims to maximize the use of solar energy, while charging as fast as possible.

## 🔋 Controller logic

The core of processing logic is contained in `internal/controller.go` module. It runs the `tick` function at set interval.
The entire flow of `tick` function is quite simple and documented - reading it will give you full perspective on how the processing works.

The following diagram presents simplified cycle flow:

```
┌───────────────┐             ┌───────────────┐        ┌ ─ ─ ─ ─ ─ ─ ─ ┐  ┌───────────────┐
│   resistere   │░            │  PV inverter  │░           Vehicle 1      │   Vehicle 2   │░
└───────────────┘░            └───────────────┘░       └ ─ ─ ─ ─ ─ ─ ─ ┘  └───────────────┘░
 ░░░░░░░│░░░░░░░░░             ░░░░░░░│░░░░░░░░░               │           ░░░░░░░│░░░░░░░░░

        │      Get current power      │                        │                  │
           production and consumption
        │───────────────────────────▶┌┴┐                       │                  │
                                     │ │
        │                            │ │                       │                  │
         ◀───────────────────────────└─┘
        │                             │                        │                  │
       ┌─┐───┐
       │ │   │Calculate energy        │                        │                  │
       │ │   │surplus
       └┬┘◀──┘                        │                        │                  │

        │                   Is vehicle in range?               │                  │
         ────────────────────────────────────────────────────▶┌─┐
        │                             │                       │ │                 │
                                                              │ │
        │                    Connection timeout               │ │                 │
         ◀────────────────────────────────────────────────────└─┘
        │                             │                        │                  │
                                         Is vehicle in range?
        │─────────────────────────────┼────────────────────────┼────────────────▶┌┴┐
                                                                                 │ │
        │                             │                        │                 │ │
                                        In range and charging                    │ │
        │◀────────────────────────────┼────────────────────────┼─────────────────└┬┘

       ┌┴┐───┐                        │                        │                  │
       │ │   │Calculate 𝚫 amps
       │ │   │                        │                        │                  │
       └─┘◀──┘
        │                             │  Set charging amps     │                  │
         ────────────────────────────────────────────────────────────────────────▶
        │                             │                        │                  │
┌───────────────┐             ┌───────────────┐        ┌ ─ ─ ─ ─ ─ ─ ─ ┐  ┌───────────────┐
│   resistere   │░            │  PV inverter  │░           Vehicle 1      │   Vehicle 2   │░
└───────────────┘░            └───────────────┘░       └ ─ ─ ─ ─ ─ ─ ─ ┘  └───────────────┘░
 ░░░░░░░░░░░░░░░░░             ░░░░░░░░░░░░░░░░░                           ░░░░░░░░░░░░░░░░░
```

To calculate delta amps, controller has to calculate current power surplus (difference between power production and consumption) and read current charging amps from vehicle.

It then calculates it using the following formula:

$$
\Delta A = \frac{W}{V \cdot 3}
$$

Where:

- $\Delta A$ - the amount by which we can increase or decrease charging current.
- $W$ - the energy surplus (the difference between power generated by PV and total power used).
- $V$ - the electric potential of the energy grid.
- $3$ - represents that the EVSE (Electric Vehicle Supply Equipment) is using three-phases to charge.

## 🚀 Deployment

### Hardware

I've done this project for myself and while it's fairly configurable, it does require specific hardware combinations to work:

- Sofar HYD 5-20KTL-3PH inverter with Wi-Fi data logger (but it should work with any inverter that outputs the Modbus data in similar way).
- Linux-based device with Bluetooth Low Energy connectivity (like Raspberry Pi Zero 2 W).
- Any Tesla vehicle that can be controlled via Bluetooth.

Adding support for other inverters should be easy, provided you can connect to them and read or calculate energy surplus.

### Software

- Download or build from source `resistere` binary for your architecture.
- Make sure your Linux environment has [`tesla-control`](https://github.com/teslamotors/vehicle-command/blob/main/cmd/tesla-control/README.md) and Python 3 installed and in path.
- Prepare configuration file.

### Configuration

The application requires `config.toml` file to run. Here's example configuration:

```toml
[web]
port = 5467 # Port on which the web UI will be running.

[controller]
cycle_interval_seconds = 10 # Time between charging current modulation (in seconds).
safety_margin_watts = 1000  # Surplus power safety margin (in watts).
grid_voltage = 230          # Grid voltage in your country/area (in volts).

[solarman_inverter]
ip = "192.168.1.99"   # IP address of inverter's data logger on your local network.
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

### Tesla pairing

`resistere` is using Bluetooth Low Energy to control charging current in Tesla vehicles. To communicate with the car it uses Tesla's official [`vehicle-command`](https://github.com/teslamotors/vehicle-command) library.

Refer to [library's documentation](https://github.com/teslamotors/vehicle-command/blob/main/README.md) for detailed instructions on pairing.

## 👆 Web UI

The application provides a simple user interface to change the controller mode (automatic/manual).

![Screenshot of user interface](resources/webui_screenshot.png)

## 🚧 Development

This project was built using Go 1.24.1. It uses [templ](https://templ.guide/) for web UI templating.

The frontend is using [HTMX](https://htmx.org/), [Surreal](https://github.com/gnat/surreal), [css-scope-inline](https://github.com/gnat/css-scope-inline) and [Phosphor icons](https://phosphoricons.com/).

To build it just execute:

```sh
make
```

When running debug build locally, you will want to put `simulator_mode = true` in `config.toml` and run:

```sh
make build_and_run
```

To make release build for Raspberry Pi (arm64-linux) run:

```sh
make release_rpi
```

To run all tests execute:

```sh
go test ./...
```

## 📈 Metrics

There are a couple of simple metrics in Prometheus notation available on `/metrics/prometheus`.

## 📜 License

This project is licensed under the [MIT license](LICENSE).
