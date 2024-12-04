package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Roles []string

func (roles *Roles) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &roles)
		return nil
	case string:
		json.Unmarshal([]byte(v), &roles)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (roles Roles) Value() (driver.Value, error) {
	return json.Marshal(roles)
}
