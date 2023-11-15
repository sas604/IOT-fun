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
	ID        int
	Name      string
	LastState string
	TopicBase string
}

type SwitchModel struct {
	DB   *sql.DB
	MQTT mqtt.Client
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

func (s *SwitchModel) SetState(id int, state string) error {

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

	fmt.Printf("Change status of %s to %s \n", sw.Name, sw.LastState)
	return nil
}
