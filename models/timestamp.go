package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *Timestamp) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)

	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return
}

func (t Timestamp) String() string { return time.Time(t).String() }

// This came from chatgpt i'm too lazy to do it properly later
// There seems to be a lot of pointless things
// @TODO: rewrite myself
func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(t) = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*(*time.Time)(t) = v
	case int64:
		*(*time.Time)(t) = time.Unix(v, 0)
	case float64:
		*(*time.Time)(t) = time.Unix(int64(v), 0)
	case []byte:
		q, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*(*time.Time)(t) = time.Unix(q, 0)
	case string:
		q, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		*(*time.Time)(t) = time.Unix(q, 0)
	default:
		return fmt.Errorf("unsupported Scan type for Timestamp: %T", value)
	}
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	return time.Time(t).Unix(), nil
}
