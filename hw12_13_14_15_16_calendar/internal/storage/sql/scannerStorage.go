package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Calendar/hw12_13_14_15_calendar/internal/kafka"
	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
)

type StorageScanner struct {
	DB     *sql.DB
	DBConf *kafka.DBConfig
}

func NewScanner(d *kafka.DBConfig, ctx context.Context) *StorageScanner {
	strg := &StorageScanner{DBConf: d}
	err := strg.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return strg
}

func (s *StorageScanner) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.DBConf.Host,
		s.DBConf.Port,
		s.DBConf.User,
		s.DBConf.Password,
		s.DBConf.Database,
		"disable",
	)
	var err error
	s.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}
	er := s.DB.PingContext(ctx)
	if er != nil {
		fmt.Println(er)
		return fmt.Errorf("failed to connect to DB: %w", er)
	}
	return nil
}
func (s *StorageScanner) DeleteNotificationTime(ctx context.Context, timeDelete time.Duration) error {
	now := time.Now().Truncate(time.Hour).Add(-timeDelete)
	_, err := s.DB.ExecContext(ctx, `
				DELETE 
				FROM calendar.events
				WHERE duration_of_the_event  <= $1
				`, now)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageScanner) GetNotificationTimeUntilEvent(ctx context.Context, checkInterval int) ([]*storage.Notification, error) {
	now := time.Now().Truncate(time.Minute)
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, title, date_and_time_of_the_event, id_user 
		FROM calendar.events 
		WHERE date_and_time_of_the_event - (time_until_event * INTERVAL '1 MINUTE') - (3 * INTERVAL '1 H') 
		    >= $1
		    AND date_and_time_of_the_event - (time_until_event * INTERVAL '1 MINUTE') - (3 * INTERVAL '1 H') < $2   `, now, now.Add(time.Minute*time.Duration(checkInterval)))
	if err != nil {
		return nil, err
	}
	var notification []*storage.Notification
	for rows.Next() {
		var not storage.Notification
		if err = rows.Scan(&not.ID, &not.Title, &not.DataStart, &not.IdUser); err != nil {
			log.Println(err)
			continue
		}
		notification = append(notification, &not)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notification, nil
}
func (s *StorageScanner) CreateNotification(ctx context.Context, notice *storage.Notification) error {
	_, err := s.DB.ExecContext(ctx, `
				insert into calendar.notification (
								   id, 
								   title,
								   date_and_time_of_the_event,
								   id_user
								   )
									values($1, $2, $3, $4)`,
		notice.ID,
		notice.Title,
		notice.DataStart,
		notice.IdUser)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageScanner) UpdateNotification(ctx context.Context, notice *storage.Notification) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE calendar.notification 
			SET title = $2,
			    date_and_time_of_the_event = $3,
			    id_user = $6
			WHERE id = $1`,
		notice.ID,
		notice.Title,
		notice.DataStart,
		notice.IdUser)
	if err != nil {
		return err
	}
	return nil
}
