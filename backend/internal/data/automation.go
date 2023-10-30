package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Automation struct {
	Abbriviation       string
	DisplayValue       string
	Switch             int
	CurrentSwitchState string
	MaxValue           float64
	MinValue           float64
}

type AutomationModel struct {
	DB *sql.DB
}

func (a *AutomationModel) GetAutomationData(abb string) (*Automation, error) {
	if abb == "" {
		return nil, ErrRecordNotFound
	}
	query := `SELECT measurements.abbreviation, measurements.display_value, switches.id, switches.state, targets.max_value, targets.min_value
	FROM measurements
	INNER JOIN switches_measurements  ON switches_measurements.measurement_id = measurements.id
	INNER JOIN switches ON switches_measurements.switch_id = switches.id
    INNER JOIN targets ON measurements.id = targets.measurement_id
	WHERE measurements.abbreviation = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var au Automation
	err := a.DB.QueryRowContext(ctx, query, abb).Scan(&au.Abbriviation, &au.DisplayValue, &au.Switch, &au.CurrentSwitchState, &au.MaxValue, &au.MinValue)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &au, nil
}
