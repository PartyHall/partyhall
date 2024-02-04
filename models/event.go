package models

type Event struct {
	Id       int       `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Date     Timestamp `json:"date" db:"date"`
	Author   string    `json:"author" db:"author"`
	Location *string   `json:"location" db:"location"`

	Exporting  bool       `json:"exporting" db:"exporting"`
	LastExport *Timestamp `json:"last_export" db:"last_export"`

	AmtImagesHandtaken  int `json:"amt_images_handtaken" db:"amt_images_handtaken"`
	AmtImagesUnattended int `json:"amt_images_unattended" db:"amt_images_unattended"`
}

type ExportedEvent struct {
	Id       int64      `json:"id" db:"id"`
	EventId  int64      `json:"event_id" db:"event_id"`
	Filename string     `json:"filename" db:"filename"`
	Date     *Timestamp `json:"date" db:"date"`
}
