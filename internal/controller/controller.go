package controller

import (
	"github.com/deseteral/resistere/internal/configuration"
	"log"
	"time"
)

func StartController(config *configuration.Controller) {
	log.Printf("Starting controller with %v second interval.\n", config.CycleIntervalSeconds)

	interval := time.Duration(config.CycleIntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)

	// Run first tick before ticker starts.
	tick()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				tick()
			}
		}
	}()
}

func tick() {
	log.Println("Starting controller tick.")
}
