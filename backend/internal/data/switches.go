package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Switch struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	LastState string `json:"state"`
	TopicBase string `json:"topic"`
}

type SwitchModel struct {
	DB   *sql.DB
	MQTT mqtt.Client
}

type SwitchMetadata struct {
	Name               string        `json:"name"`
	State              string        `json:"state"`
	Value              float64       `json:"value"`
	MeasurementDisplay string        `json:"unit"`
	MeasurementAbb     string        `json:"abbriviation"`
	Measurement        string        `json:"measurement"`
	Automation         bool          `json:"automation"`
	MaxValue           float64       `json:"maxValue"`
	MinValue           float64       `json:"minValue"`
	Schedule           bool          `json:"schedule"`
	Interval           time.Duration `json:"interval"`
	Duration           time.Duration `json:"duration"`
}

func (s SwitchModel) Insert(sw *Switch) error {
	query := `
		INSERT INTO switches (name, state)
		VALUES ($1, $2)
`
	args := []any{sw.Name, sw.LastState}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&sw.ID)

}

func (s SwitchModel) GetOneByID(id int) (*Switch, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, name, state, topic_base
	FROM switches
	WHERE id = $1`

	var sw Switch

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&sw.ID, &sw.Name, &sw.LastState, &sw.TopicBase)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &sw, nil
}

func (s SwitchModel) GetAll() ([]*Switch, error) {
	query := `
	SELECT id, name, state, topic_base
	FROM switches`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	rows, err := s.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	switches := []*Switch{}
	for rows.Next() {
		var sw Switch
		err := rows.Scan(&sw.ID, &sw.Name, &sw.LastState, &sw.TopicBase)

		if err != nil {
			return nil, err
		}
		switches = append(switches, &sw)
	}
	return switches, nil
}

func (s *SwitchModel) SetState(id int, state string) error {
	fmt.Printf("Switch %d set to %s \n", id, state)
	sw, err := s.GetOneByID(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	query := `
		UPDATE switches
		SET state = $1
		WHERE id = $2
		RETURNING state`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, query, state, sw.ID).Scan(&sw.LastState)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	s.MQTT.Publish(sw.TopicBase+"/set/"+strconv.Itoa(sw.ID), 0, true, sw.LastState)
	return nil
}

func (s *SwitchModel) GetAllWithMetaData(measurmentsMap map[string]float64) ([]*SwitchMetadata, error) {
	data := []*SwitchMetadata{}
	query := `
	SELECT switches.name, switches.state, 
	CASE WHEN measurements.abbreviation IS NULL THEN  "false" ELSE "true" END AS automation,
	
	ifnull(measurements.abbreviation, '') AS abbreviation, ifnull(measurements.display_value,'') AS display_value, 
	
	CASE WHEN jobs.id IS NULL THEN  "false" ELSE "true" END AS schedule,
	
	ifnull(targets.max_value, 0) AS max_value, ifnull(targets.min_value, 0) AS min_value,  ifnull(targets.display_value, '') AS display_value,  ifnull(jobs.interval, 0) AS interval,  ifnull(jobs.duration, 0) AS duration

	FROM switches
	LEFT JOIN switches_measurements ON switches_measurements.switch_id = switches.id
	LEFT JOIN measurements ON switches_measurements.measurement_id = measurements.id
	LEFT JOIN targets ON measurements.id = targets.measurement_id
	LEFT JOIN jobs ON jobs.switch = switches.id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var s SwitchMetadata
		err = rows.Scan(&s.Name, &s.State, &s.Automation, &s.MeasurementAbb, &s.Measurement, &s.Schedule, &s.MaxValue, &s.MinValue, &s.MeasurementDisplay, &s.Interval, &s.Duration)
		if err != nil {
			return nil, err
		}
		s.Value = measurmentsMap[s.MeasurementAbb]
		data = append(data, &s)
	}

	return data, nil
}
