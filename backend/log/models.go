package log

import "time"

type Log struct {
	Id        int       `db:"id"`
	Type      string    `db:"type" json:"type"`
	Text      string    `db:"text" json:"text"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}
