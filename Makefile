all: solarman_interface_binary templ build run

build:
	@echo "Building application"
	go build -o ./bin/resistere ./cmd/resistere
	@echo ""

build_release:
	@echo "Building release version of application"
	go build -ldflags "-s -w" -o ./bin/resistere ./cmd/resistere
	@echo ""

run:
	@echo "Running application"
	./bin/resistere

templ:
	@echo "Building templ components"
	go tool templ generate
	@echo ""

solarman_interface_binary:
	@echo "Building solarman_interface Python binary"
	mkdir -p ./internal/inverter/solarman_interface/build
	python3 -m zipapp ./internal/inverter/solarman_interface/src -m "main:main" -o ./internal/inverter/solarman_interface/build/solarman_interface.pyz -p "/usr/bin/env python3"
	chmod u+x ./internal/inverter/solarman_interface/build/solarman_interface.pyz
	@echo ""
