package controller

import (
	"errors"
	"fmt"
)

type Mode int

const (
	ModePVAutomatic Mode = iota
	ModeManual
	ModeManualUntilDisconnect
)

var modeName = map[Mode]string{
	ModePVAutomatic:           "PV Automatic",
	ModeManual:                "Manual",
	ModeManualUntilDisconnect: "Manual until unplugged",
}

func ParseIntToMode(value int) (Mode, error) {
	if value < int(ModePVAutomatic) || value > int(ModeManualUntilDisconnect) {
		return -1, errors.New(fmt.Sprintf("could not parse value %d to Mode", value))
	}

	return Mode(value), nil
}
