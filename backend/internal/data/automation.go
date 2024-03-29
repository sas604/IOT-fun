package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Automation struct {
	Abbriviation       string  `json:"abbriviation"`
	DisplayValue       string  `json:"displayValue"`
	Switch             int     `json:"switch"`
	CurrentSwitchState string  `json:"currentState"`
	MaxValue           float64 `json:"max"`
	MinValue           float64 `json:"min"`
}

type Job struct {
	Interval int
	Duration int
	Switch   int
	OnStart  string
	OnEnd    string
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

func (a *AutomationModel) GetJobData() ([]*Job, error) {
	jobs := []*Job{}

	query :=
		`SELECT interval, duration, switch, on_start, on_end 
			FROM jobs`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var job Job

		err := rows.Scan(
			&job.Interval,
			&job.Duration,
			&job.Switch,
			&job.OnStart,
			&job.OnEnd,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}
