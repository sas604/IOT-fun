package main

import (
	"fmt"
	"time"
)

// check if measurement exsit and has switch and target valuse assosiated with it
// check if action needed and handle it if so.
func (app *application) handleMonitoring(measurements map[string]float64) {
	fmt.Println(measurements)
	for abb, val := range measurements {
		au, err := app.models.Automations.GetAutomationData(abb)
		fmt.Println(au)
		fmt.Println(val)
		if err != nil {
			app.logger.Error("error gettingautomations", "error", err.Error())
			continue
		}

		if au.MinValue > val {
			err = app.models.Switches.SetState(au.Switch, "on")
			if err != nil {
				app.logger.Error(err.Error())
			}

		}
		if au.MaxValue < val {
			err = app.models.Switches.SetState(au.Switch, "off")
			if err != nil {
				app.logger.Error(err.Error())
			}

		}

	}

}

func (app *application) handlePeriodicTasks(interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

}
