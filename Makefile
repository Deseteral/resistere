all: solarman_interface_binary templ build

build:
	@echo "Building application"
	go build -o ./bin/resistere ./cmd/resistere
	@echo ""

run:
	@echo "Running application"
	./bin/resistere

build_and_run: build run

templ:
	@echo "Building templ components"
	go tool templ generate
	@echo ""

solarman_interface_binary:
	@echo "Building solarman_interface Python binary"
	mkdir -p ./internal/pv/solarman_interface/build
	python3 -m zipapp ./internal/pv/solarman_interface/src -m "main:main" -o ./internal/pv/solarman_interface/build/solarman_interface.pyz -p "/usr/bin/env python3"
	chmod u+x ./internal/pv/solarman_interface/build/solarman_interface.pyz
	@echo ""

release_rpi: templ solarman_interface_binary
	mkdir -p "./bin/release"

	@echo "Building release version of application"
	GOOS="linux" GOARCH="arm64" go build -ldflags "-s -w" -o ./bin/release/resistere ./cmd/resistere
	@echo ""

	@echo "Copy scripts"
	cp ./scripts/* ./bin/release/
	@echo ""

	@echo "Make zip bundle"
	zip -vrj ./bin/resistere-arm64-linux.zip ./bin/release/ -X
