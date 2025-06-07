package metrics

import "time"

type Registry struct {
	LatestFrame Frame
}

type Frame struct {
	Timestamp             time.Time
	PowerProductionWatts  float64
	PowerConsumptionWatts float64
	VehicleFrames         []VehicleFrame
}

type VehicleFrame struct {
	Name               string
	ChargingPowerWatts float64
}

func NewMetricsRegistry() *Registry {
	return &Registry{}
}

func NewMetricsFrame() Frame {
	return Frame{
		Timestamp: time.Now(),
	}
}

func NewMetricsVehicleFrame(name string) VehicleFrame {
	return VehicleFrame{
		Name:               name,
		ChargingPowerWatts: 0.0,
	}
}
