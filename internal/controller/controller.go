package controller

import (
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/pv"
	"log"
	"time"
)

type Controller struct {
	updateInterval time.Duration
	inverter       pv.Inverter
}

func (c Controller) StartBackgroundTask() {
	log.Printf("Starting controller with %v interval.\n", c.updateInterval)

	// Run first tick before ticker starts.
	c.tick()

	ticker := time.NewTicker(c.updateInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.tick()
			}
		}
	}()
}

func (c Controller) tick() {
	log.Println("Starting controller tick.")
}

func NewController(inverter pv.Inverter, config *configuration.Controller) Controller {
	return Controller{
		updateInterval: time.Duration(config.CycleIntervalSeconds) * time.Second,
		inverter:       inverter,
	}
}
