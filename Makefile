all: templ build run

build:
	@echo "Building application"
	go build -o ./bin/resistere ./cmd/resistere
	@echo ""

run:
	@echo "Running application"
	./bin/resistere

templ:
	@echo "Building templ components"
	~/go/bin/templ generate
	@echo ""

