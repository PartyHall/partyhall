package models

type AppState struct {
	HardwareID            string  `json:"hwid" db:"hwid"`
	ApiToken              *string `json:"token" db:"token"`
	CurrentEvent          *int    `json:"-" db:"current_event"`
	LastAppliedMigrations int     `json:"-" db:"last_applied_migration"`

	CurrentEventObj *Event `json:"current_event" db:"-"`
}
