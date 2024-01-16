package main

import (
	"net/http"
)

func (app *application) listMeasurementsHandler(w http.ResponseWriter, r *http.Request) {

	ms, err := app.models.Measurements.GetMeasurementMap()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	sw, err := app.models.Switches.GetAllWithMetaData(ms)

	data := envelope{
		"measurments": ms,
		"switches":    sw,
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJSON(w, 200, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
