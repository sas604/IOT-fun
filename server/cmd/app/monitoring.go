package main

import (
	"fmt"
)

// check if measurement exsit and has switch and target valuse assosiated with it
// check if action needed and handle it if so.
func (app *application) handleMonitoring(measurements map[string]float64) {

	for abb, val := range measurements {
		au, err := app.models.Automations.GetAutomationData(abb)
		if err != nil {
			app.logger.Error("error gettingautomations", "error", err.Error())
			continue
		}
		fmt.Println(abb, au.Abbriviation, au.MinValue, au.MaxValue, au.CurrentSwitchState, val)
		switch au.Abbriviation {
		case "temp":
			if au.MinValue > val && au.CurrentSwitchState == "off" {
				err = app.models.Switches.SetState(au.Switch, "on")
				app.logger.Error(err.Error())
			}
			if au.MaxValue < val && au.CurrentSwitchState == "on" {
				app.models.Switches.SetState(au.Switch, "off")
				app.logger.Error(err.Error())
			}
		}

	}

}
