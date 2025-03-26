package dal

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/state"
)

var EVENTS Events

type Events struct{}

func (e Events) GetCollection(amt, offset int) (*models.PaginatedResponse, error) {
	resp := models.PaginatedResponse{}

	row := DB.QueryRow(`SELECT COUNT(*) FROM event`)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&resp.TotalCount)
	if err != nil {
		return nil, err
	}

	resp.CalculateMaxPage()

	rows, err := DB.Queryx(`
		SELECT
			id,
			name,
			author,
			date,
			location,
			nexus_id,
			amt_images_handtaken,
			amt_images_unattended
		FROM event
		LIMIT ?
		OFFSET ?
	`, amt, offset)

	if err != nil {
		return nil, err
	}

	events := []models.Event{}

	for rows.Next() {
		event := models.Event{}

		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	resp.Results = events

	return &resp, nil
}

func (e Events) GetCurrent() (*models.Event, error) {
	row := DB.QueryRowx(`
		SELECT
			e.id AS id,
			e.name AS name,
			e.author AS author,
			e.date AS date,
			e.location AS location,
			e.nexus_id AS nexus_id,
			e.amt_images_handtaken AS amt_images_handtaken,
			e.amt_images_unattended AS amt_images_unattended
		FROM event e
		INNER JOIN app_state app ON e.id = app.current_event
		LIMIT 1
	`)

	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, nil
		}

		return nil, row.Err()
	}

	event := models.Event{}
	err := row.StructScan(&event)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &event, err
}

func (e Events) Create(evt *models.Event) error {
	_, err := DB.Exec(`
		INSERT INTO event (name, author, date, location, nexus_id)
		VALUES (?, ?, ?, ?, ?);
	`, evt.Name, evt.Author, evt.Date, evt.Location, evt.NexusId)

	if err != nil {
		return err
	}

	row := DB.QueryRow(`SELECT id FROM event WHERE rowid = last_insert_rowid();`)

	if row.Err() != nil {
		return row.Err()
	}

	var id int64
	err = row.Scan(&id)

	if err != nil {
		return err
	}

	evt.Id = id

	return nil
}

func (e Events) Update(evt *models.Event) error {
	_, err := DB.Exec(`
		UPDATE event
		SET name = ?, author = ?, date = ?, location = ?, nexus_id = ?
		WHERE id = ?;
	`, evt.Name, evt.Author, evt.Date, evt.Location, evt.NexusId, evt.Id)

	if err != nil {
		return err
	}

	return nil
}

func (e Events) Delete(id int64) error {
	_, err := DB.Exec(`DELETE FROM event WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

func (e Events) Get(id int64) (*models.Event, error) {
	row := DB.QueryRowx(`
		SELECT
			id AS id,
			name AS name,
			author AS author,
			date AS date,
			location AS location,
			nexus_id AS nexus_id,
			amt_images_handtaken AS amt_images_handtaken,
			amt_images_unattended AS amt_images_unattended
		FROM event
		WHERE id = ?
	`, id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var event models.Event
	err := row.StructScan(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e Events) GetAndSetAny() (*models.Event, error) {
	row := DB.QueryRowx(`
		SELECT
			e.id AS id,
			e.name AS name,
			e.author AS author,
			e.date AS date,
			e.location AS location,
			e.nexus_id AS nexus_id,
			e.amt_images_handtaken AS amt_images_handtaken,
			e.amt_images_unattended AS amt_images_unattended
		FROM event e
		LIMIT 1
	`)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var event models.Event

	err := row.StructScan(&event)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &event, e.Set(&event)
}

func (e Events) Set(event *models.Event) error {
	_, err := DB.Exec(`
		UPDATE app_state
		SET current_event = ?
	`, event.Id)

	return err
}

func (e Events) InsertPicture(eventId int64, filename string, unattended bool, alternateFilename string) (string, error) {
	id := uuid.NewString()

	_, err := DB.Exec(`
		INSERT INTO picture(id, taken_at, unattended, filename, event_id, nexus_id, alternate_filename)
		VALUES (?, ?, ?, ?, ?, NULL, ?);
	`, id, time.Now(), unattended, filename, state.STATE.CurrentEvent.Id, alternateFilename)

	if err != nil {
		return "", err
	}

	// Re-calculating amount of images
	_, err = DB.Exec(`
		UPDATE event SET amt_images_unattended  = (
			SELECT COUNT(*)
			FROM event e
			INNER JOIN app_state app ON app.current_event  = e.id
			INNER JOIN picture p ON p.event_id  = e.id
			WHERE p.unattended IS TRUE
		), amt_images_handtaken  = (
			SELECT COUNT(*)
			FROM event e
			INNER JOIN app_state app ON app.current_event  = e.id
			INNER JOIN picture p ON p.event_id  = e.id
			WHERE p.unattended IS FALSE
		)
		WHERE event.id  = (SELECT current_event FROM app_state LIMIT 1)
	`)

	if err != nil {
		log.Error("Failed to update picture count", "err", err)
	}

	return id, nil
}

func (e Events) GetPictures(eventId int64, isUnattended bool) ([]models.Picture, error) {
	pictures := []models.Picture{}

	rows, err := DB.Queryx(`
		SELECT id, taken_at, unattended, filename, event_id, nexus_id
		FROM picture
		WHERE event_id = ? AND unattended = ?
	`, eventId, isUnattended)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := models.Picture{}

		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}

		pictures = append(pictures, p)
	}

	return pictures, nil
}

func (e Events) GetUnsubmittedPictures(eventId int64) ([]models.Picture, error) {
	pictures := []models.Picture{}

	rows, err := DB.Queryx(`
		SELECT id, taken_at, unattended, filename, event_id, nexus_id
		FROM picture
		WHERE event_id = ? AND nexus_id IS NULL
	`, eventId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := models.Picture{}

		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}

		pictures = append(pictures, p)
	}

	return pictures, nil
}

func (e Events) UpdatePicture(p models.Picture) error {
	_, err := DB.Exec(`
		UPDATE picture
		SET taken_at = ?, unattended = ?, filename = ?, nexus_id = ?
		WHERE id = ?
	`, p.TakenAt, p.Unattended, p.Filename, p.NexusId, p.Id)

	return err
}
