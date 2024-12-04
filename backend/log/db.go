package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/partyhall/partyhall/mercure_client"
	"go.uber.org/zap/zapcore"
)

func Message(level zapcore.Level, msg string, args []any) error {
	fullTextArr := []string{msg}
	data := []any{}

	for _, r := range args {
		data = append(data, r)
		fullTextArr = append(fullTextArr, fmt.Sprintf("%v", r))
	}

	fullText := strings.Join(fullTextArr, " ")

	LOG.Logw(level, msg, data...)

	if DB != nil {
		go func() {
			ts := time.Now()
			_, err := DB.Exec(`
				INSERT INTO logs (type, text, timestamp)
				VALUES (?, ?, ?);
			`, level.CapitalString(), fullText, ts)

			if err != nil {
				LOG.Errorw("Failed to write log in db", "err", err)
				return
			}

			row := DB.QueryRow(`SELECT id FROM logs WHERE rowid = last_insert_rowid();`)
			if row.Err() == nil {
				var id int
				row.Scan(&id)

				mercure_client.CLIENT.PublishEvent("/logs", Log{
					Id:        id,
					Type:      level.CapitalString(),
					Text:      fullText,
					Timestamp: ts,
				})

				return
			}

			LOG.Errorw("Failed to write log in db", "err", row.Err())
		}()
	}

	return nil
}

func Error(msg string, args ...any) error {
	return Message(zapcore.ErrorLevel, msg, args)
}

func Warn(msg string, args ...any) error {
	return Message(zapcore.WarnLevel, msg, args)
}

func Info(msg string, args ...any) error {
	return Message(zapcore.InfoLevel, msg, args)
}

func Debug(msg string, args ...any) error {
	return Message(zapcore.DebugLevel, msg, args)
}

func CountMessages() (int, error) {
	return 0, nil
}

/** @TODO: No limit/offset but limit/lastLogGotten **/
func GetMessages(limit, offset int) ([]Log, error) {
	rows, err := DB.Queryx(`
		SELECT id, type, text, timestamp
		FROM logs
		ORDER BY id DESC
		LIMIT ?
		OFFSET ?
	`, limit, offset)

	if err != nil {
		return nil, err
	}

	logs := []Log{}
	for rows.Next() {
		var log Log

		err = rows.StructScan(&log)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}
