package vehicle

import (
	"context"
	"fmt"
	"time"

	"github.com/deseteral/resistere/internal/configuration"
	"github.com/teslamotors/vehicle-command/pkg/connector/ble"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
	"github.com/teslamotors/vehicle-command/pkg/protocol/protobuf/carserver"
	"github.com/teslamotors/vehicle-command/pkg/protocol/protobuf/universalmessage"
	tesla "github.com/teslamotors/vehicle-command/pkg/vehicle"
)

type chargingVehicle interface {
	GetState(context.Context, tesla.StateCategory) (*carserver.VehicleData, error)
	SetChargingAmps(context.Context, int32) error
	Disconnect()
}

type TeslaControlController struct {
	connect func(vin string) (chargingVehicle, error)
}

func (c *TeslaControlController) GetChargingState(vehicle *Vehicle) (*ChargingState, error) {
	car, err := c.connect(vehicle.Vin)
	if err != nil {
		return nil, err
	}
	defer car.Disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := car.GetState(ctx, tesla.StateCategoryCharge)
	if err != nil {
		return nil, fmt.Errorf("get charging state: %w", err)
	}
	chargeState := data.GetChargeState()
	if chargeState == nil || chargeState.GetChargingState() == nil {
		return nil, fmt.Errorf("vehicle response is missing charging state")
	}
	if chargeState.GetChargingState().GetCharging() == nil {
		return nil, nil
	}
	if chargeState.GetOptionalChargingAmps() == nil || chargeState.GetOptionalChargerPower() == nil {
		return nil, fmt.Errorf("vehicle response is missing charging amps or power")
	}

	return &ChargingState{
		Amps:  int(chargeState.GetChargingAmps()),
		Power: int(chargeState.GetChargerPower()),
	}, nil
}

func (c *TeslaControlController) SetChargingAmps(vehicle *Vehicle, amps int) error {
	car, err := c.connect(vehicle.Vin)
	if err != nil {
		return err
	}
	defer car.Disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := car.SetChargingAmps(ctx, int32(amps)); err != nil {
		return fmt.Errorf("set charging amps: %w", err)
	}
	return nil
}

func NewTeslaControlController(config *configuration.TeslaControl) (*TeslaControlController, error) {
	key, err := protocol.LoadPrivateKey(config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load Tesla private key: %w", err)
	}

	return &TeslaControlController{
		connect: func(vin string) (chargingVehicle, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			conn, err := ble.NewConnection(ctx, vin)
			if err != nil {
				return nil, fmt.Errorf("connect to Tesla over BLE: %w", err)
			}
			car, err := tesla.NewVehicle(conn, key, nil)
			if err != nil {
				conn.Close()
				return nil, fmt.Errorf("create Tesla vehicle: %w", err)
			}
			if err := car.Connect(ctx); err != nil {
				car.Disconnect()
				return nil, fmt.Errorf("connect to Tesla vehicle: %w", err)
			}
			// Charge-state reads and charging commands only use infotainment.
			if err := car.StartSession(ctx, []universalmessage.Domain{universalmessage.Domain_DOMAIN_INFOTAINMENT}); err != nil {
				car.Disconnect()
				return nil, fmt.Errorf("authenticate Tesla session: %w", err)
			}
			return car, nil
		},
	}, nil
}
