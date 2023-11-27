package main

import (
	"context"
	"sync"
	"time"

	"github.com/sas604/IOT-fun/server/internal/data"
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

func (app *application) handlePeriodicTasks() {
	jobs, err := app.models.Automations.GetJobData()
	if err != nil {
		app.logger.Error(err.Error())
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, j := range jobs {
		wg.Add(1)
		go func(j *data.Job) {
			defer wg.Done()
			ticker := time.NewTicker(time.Duration(j.Interval) * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					ticker.Stop()
					err := app.models.Switches.SetState(j.Switch, j.OnStart)
					if err != nil {
						app.logger.Error("error in job scheduler ", "error", err)
						cancel()
						return
					}
					t := time.NewTimer(time.Duration(j.Duration) * time.Minute)
					go func() {
						<-t.C
						app.models.Switches.SetState(j.Switch, j.OnEnd)
						ticker.Reset(3 * time.Second)
					}()
				}

			}
		}(j)
	}
	wg.Wait()
}
