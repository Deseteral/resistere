package controller

import (
	"errors"
	"fmt"
)

type Mode int

const (
	ModePVAutomatic Mode = iota
	ModeManual
)

var modeName = map[Mode]string{
	ModePVAutomatic: "PV Automatic",
	ModeManual:      "Manual",
}

func ParseIntToMode(value int) (Mode, error) {
	if value > int(ModeManual) {
		return -1, errors.New(fmt.Sprintf("could not parse value %d to Mode", value))
	}

	return Mode(value), nil
}
