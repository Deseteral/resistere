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
	SetChargingAmps    int
	ChargingPowerWatts float64
	IsInRange          bool
	IsSelected         bool
}

func (v *VehicleFrame) IsCharging() bool {
	return v.ChargingPowerWatts > 0
}

func NewMetricsRegistry() *Registry {
	return &Registry{}
}

func NewVehicleMetricsFrame() VehicleFrame {
	return VehicleFrame{
		SetChargingAmps:    0,
		ChargingPowerWatts: 0,
		IsInRange:          false,
		IsSelected:         false,
	}
}
