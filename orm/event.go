package orm

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/partyhall/partyhall/models"
)

type Events struct {
	db *sqlx.DB
}

func (e *Events) ClearExporting() error {
	_, err := e.db.Exec(`
		UPDATE event
		SET exporting = FALSE
	`)

	return err
}

func (e *Events) GetEvents() ([]models.Event, error) {
	events := []models.Event{}

	// @TODO created_at && order by created_at desc
	// @TODO rewrite with no subquery
	rows, err := e.db.Queryx(`
		SELECT id, name, date, author, location, exporting, last_export, COALESCE(counts.handtaken, 0) amt_images_handtaken, COALESCE(counts.unattended, 0) amt_images_unattended
		FROM event
		LEFT JOIN (
			SELECT event_id,
				SUM(CASE unattended WHEN TRUE THEN 0 ELSE 1 END) handtaken,
				SUM(CASE unattended WHEN TRUE THEN 1 ELSE 0 END) unattended
			FROM image
			GROUP BY event_id
		) counts ON event.id = counts.event_id
		ORDER BY name
	`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		evt := models.Event{}
		err = rows.StructScan(&evt)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	return events, nil
}

func (e *Events) GetEvent(id int64) (*models.Event, error) {
	row := e.db.QueryRowx(`
		SELECT id, name, date, author, location, exporting, last_export, COALESCE(counts.handtaken, 0) amt_images_handtaken, COALESCE(counts.unattended, 0) amt_images_unattended
		FROM event
		LEFT JOIN (
			SELECT event_id,
				SUM(CASE unattended WHEN TRUE THEN 0 ELSE 1 END) handtaken,
				SUM(CASE unattended WHEN TRUE THEN 1 ELSE 0 END) unattended
			FROM image
			GROUP BY event_id
		) counts ON event.id = counts.event_id
		WHERE id = ?
	`, id)

	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}

		return nil, row.Err()
	}

	evt := models.Event{}
	err := row.StructScan(&evt)

	return &evt, err
}

func (e *Events) CreateEvent(name string, author, location *string, date *int64) (*models.Event, error) {
	row := e.db.QueryRowx(`
		INSERT INTO event (name, date, author, location)
		VALUES (?, ?, ?, ?)
		RETURNING *
	`, name, date, author, location)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var event models.Event = models.Event{}
	err := row.StructScan(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *Events) Save(event *models.Event) error {
	var date int64 = 0
	if event.Date != nil {
		date = (time.Time(*event.Date)).Unix()
	}

	var last_export *int64 = nil
	if event.LastExport != nil {
		exp := (time.Time(*event.LastExport)).Unix()
		last_export = &exp
	}

	// @TODO: find how to make date work with namedexec without this hack
	_, err := e.db.NamedExec(`
		UPDATE event
		SET name = :name, date = :date, author = :author, location = :location, exporting = :exporting, last_export = :last_export
		WHERE id = :id
	`, map[string]interface{}{
		"id":          event.Id,
		"name":        event.Name,
		"date":        date,
		"author":      event.Author,
		"location":    event.Location,
		"exporting":   event.Exporting,
		"last_export": last_export,
	})

	return err
}

func (e *Events) GetImages(event *models.Event, unattended bool) ([]models.Image, error) {
	images := []models.Image{}

	rows, err := e.db.Queryx(`
		SELECT id, created_at, unattended, event_id
		FROM image
		WHERE event_id = ? AND unattended = ?
	`, event.Id, unattended)

	if err != nil {
		return images, err
	}

	for rows.Next() {
		image := models.Image{}
		err := rows.StructScan(&image)
		if err != nil {
			return images, err
		}

		images = append(images, image)
	}

	return images, nil
}

func (e *Events) GetImage(id int64) (*models.Image, error) {
	row := e.db.QueryRowx(`
		SELECT id, created_at, unattended, event_id
		FROM image
		WHERE id = ?
	`, id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	pct := models.Image{}
	err := row.StructScan(&pct)

	return &pct, err
}

func (e *Events) InsertImage(eventId int64, unattended bool) (*models.Image, error) {
	currTime := time.Now().Unix()

	row := e.db.QueryRowx(`
		INSERT INTO image(event_id, unattended, created_at)
		VALUES (?, ?, ?)
		RETURNING *
	`, eventId, unattended, currTime)
	if row.Err() != nil {
		return nil, row.Err()
	}

	img := models.Image{}
	err := row.StructScan(&img)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func (e *Events) InsertExportedEvent(event *models.Event, filename string) (*models.ExportedEvent, error) {
	currTime := time.Now().Unix()

	row := e.db.QueryRowx(`
		INSERT INTO exported_event(event_id, filename, date)
		VALUES (?, ?, ?)
		RETURNING *
	`, event.Id, filename, currTime)
	if row.Err() != nil {
		return nil, row.Err()
	}

	exportedEvent := models.ExportedEvent{}
	err := row.StructScan(&exportedEvent)
	if err != nil {
		return nil, err
	}

	return &exportedEvent, nil
}

func (e *Events) GetExportedEvent(exportedEventId int64) (*models.ExportedEvent, error) {
	row := e.db.QueryRowx(`
		SELECT id, event_id, filename, date
		FROM exported_event
		WHERE id = ?
	`, exportedEventId)

	if row.Err() != nil {
		return nil, row.Err()
	}

	exportedEvent := models.ExportedEvent{}
	err := row.StructScan(&exportedEvent)

	return &exportedEvent, err
}

func (e *Events) GetExportedEvents(event *models.Event, limit int64) ([]models.ExportedEvent, error) {
	exportedEvents := []models.ExportedEvent{}

	rows, err := e.db.Queryx(`
		SELECT id, event_id, filename, date
		FROM exported_event
		WHERE event_id = ?
		ORDER BY date desc
		LIMIT ?
	`, event.Id, limit)

	if err != nil {
		return exportedEvents, err
	}

	for rows.Next() {
		exportedEvent := models.ExportedEvent{}
		err := rows.StructScan(&exportedEvent)
		if err != nil {
			return exportedEvents, err
		}

		exportedEvents = append(exportedEvents, exportedEvent)
	}

	return exportedEvents, nil
}
