package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type JsonnableNullstring sql.NullString

func (s JsonnableNullstring) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}

	return []byte(`null`), nil
}

func (s *JsonnableNullstring) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.String, s.Valid = "", false

		return nil
	}

	err := json.Unmarshal(data, &s.String)
	s.Valid = (err == nil)

	return err
}

func (s JsonnableNullstring) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}

	return s.String, nil
}

func (s *JsonnableNullstring) Scan(value interface{}) error {
	if value == nil {
		*s = JsonnableNullstring{
			String: "",
			Valid:  false,
		}

		return nil
	}

	switch v := value.(type) {
	case string:
		*s = JsonnableNullstring{
			String: v,
			Valid:  true,
		}
	case []byte:
		*s = JsonnableNullstring{
			String: string(v),
			Valid:  true,
		}
	default:
		return fmt.Errorf("cannot scan type %T into JsonnableNullstring", value)
	}

	return nil
}

type JsonnableNullTime sql.NullTime

func (t JsonnableNullTime) String() string {
	if !t.Valid {
		return "JNTime[value=null]"
	}

	return fmt.Sprintf("JNTime[value=%v]", t.Time.Format("2006-01-02@15:04:05"))
}

func (t *JsonnableNullTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Time)
	}
	return []byte("null"), nil
}

func (t *JsonnableNullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	err := json.Unmarshal(data, &t.Time)
	t.Valid = (err == nil)
	return err
}

func (t JsonnableNullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JsonnableNullTime) Scan(value interface{}) error {
	if value == nil {
		*t = JsonnableNullTime{
			Time:  time.Time{},
			Valid: false,
		}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = JsonnableNullTime{
			Time:  v,
			Valid: true,
		}
	default:
		return fmt.Errorf("cannot scan type %T into JsonnableNullTime", value)
	}
	return nil
}

type JsonnableNullInt64 sql.NullInt64

func (i *JsonnableNullInt64) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int64)
	}
	return []byte("null"), nil
}

func (i *JsonnableNullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		i.Int64, i.Valid = 0, false
		return nil
	}
	err := json.Unmarshal(data, &i.Int64)
	i.Valid = (err == nil)
	return err
}

func (i JsonnableNullInt64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

func (i *JsonnableNullInt64) Scan(value interface{}) error {
	if value == nil {
		*i = JsonnableNullInt64{
			Int64: 0,
			Valid: false,
		}
		return nil
	}

	switch v := value.(type) {
	case int64:
		*i = JsonnableNullInt64{
			Int64: v,
			Valid: true,
		}
	case []byte:
		return json.Unmarshal(v, &i.Int64)
	default:
		return fmt.Errorf("cannot scan type %T into JsonnableNullInt64", value)
	}
	return nil
}
