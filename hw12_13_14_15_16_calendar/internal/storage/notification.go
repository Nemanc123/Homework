package storage

import "time"

type Notification struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	DataStart time.Time `json:"date_and_time_of_the_event"`
	IdUser    int       `json:"id_user"`
}
