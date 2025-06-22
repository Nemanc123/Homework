package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"log"
)

type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Schema   string `yaml:"schema"`
	Database string `yaml:"database_db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Storage struct {
	db     *sql.DB
	dbConf DatabaseConf
	ctx    context.Context
}

func New(d DatabaseConf, ctx context.Context) *Storage {
	strg := &Storage{dbConf: d, ctx: ctx}
	err := strg.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return strg
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.dbConf.Host,
		s.dbConf.Port,
		s.dbConf.User,
		s.dbConf.Password,
		s.dbConf.Database,
		"disable",
	)
	var err error
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}
	er := s.db.PingContext(ctx)
	if er != nil {
		fmt.Println(er)
		return fmt.Errorf("failed to connect to db: %w", er)
	}
	s.ctx = ctx
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	tx, err := s.db.BeginTx(s.ctx, nil) // *sql.Tx
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	query := `insert into calendar.events (
                   id, 
                   title,
                   date_and_time_of_the_event,
                   duration_of_the_event,
                   description_event,
                   id_user,
                   time_until_event)
 					values($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(s.ctx, query,
		event.ID,
		event.Title,
		event.DataStart,
		event.DataEnd,
		event.Description,
		event.IdUser,
		event.TimeUntilEvent,
	)
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit to db: %w", err)
	}
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `DELETE FROM postgres_db.calendar.events WHERE id = $1`

	_, err = tx.ExecContext(s.ctx, query, id)
	if err != nil {
		return fmt.Errorf("id not found in db: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit to db: %w", err)
	}
	return nil
}

func (s *Storage) GetEvents() ([]storage.Event, error) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
        SELECT 	   id, 
                   title,
                   date_and_time_of_the_event,
                   duration_of_the_event,
                   description_event,
                   id_user,
                   time_until_event
        FROM calendar.events
    `

	rows, wrong := tx.QueryContext(s.ctx, query)
	if wrong != nil {
		return nil, fmt.Errorf("unable to get events: %w", err)
	}
	defer rows.Close()

	var tables []storage.Event
	for rows.Next() {
		var table storage.Event
		if err = rows.Scan(
			&table.ID,
			&table.Title,
			&table.DataStart,
			&table.DataEnd,
			&table.Description,
			&table.IdUser,
			&table.TimeUntilEvent,
		); err != nil {
			return nil, fmt.Errorf("unable to get event: %w", err)
		}

		tables = append(tables, table)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit to db: %w", err)
	}
	return tables, nil
}

func (s *Storage) UpdateEvent(id int, event storage.Event) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query :=
		`UPDATE postgres_db.calendar.events 
			SET title = $2,
			    date_and_time_of_the_event = $3,
			    duration_of_the_event = $4,
			    description_event = $5,
			    id_user = $6,
			    time_until_event = $7
			WHERE id = $1`

	_, wrong := tx.ExecContext(s.ctx, query,
		id,
		event.Title,
		event.DataStart,
		event.DataEnd,
		event.Description,
		event.IdUser,
		event.TimeUntilEvent,
	)
	if wrong != nil {
		return fmt.Errorf("not all parameters are specified: %w", wrong)
	}
	err = tx.Commit()
	if err != nil {
		fmt.Errorf("failed to commit to db: %w", err)
	}
	return nil

}
