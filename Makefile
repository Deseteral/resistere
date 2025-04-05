all: solarman_interface_binary templ build run

build:
	@echo "Building application"
	go build -o ./bin/resistere ./cmd/resistere
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
	mkdir -p ./solarman_interface/build
	python3 -m zipapp ./solarman_interface/src -m "main:main" -o ./solarman_interface/build/solarman_interface.pyz -p "/usr/bin/env python3"
	chmod u+x ./solarman_interface/build/solarman_interface.pyz
	@echo ""
