package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Switch struct {
	ID        int
	Name      string
	LastState string
}

type SwitchModel struct {
	DB *sql.DB
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

func (s SwitchModel) GetOneBySlug(id int64) (*Switch, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, name, state
	FROM switches
	WHERE id = $1`

	var sw Switch

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&sw.ID, &sw.Name, &sw.Name)

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
