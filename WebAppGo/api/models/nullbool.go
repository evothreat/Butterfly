package models

import (
	"database/sql/driver"
	"encoding/json"
)

// NullBool represents a bool that may be null.
// NullBool implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.
func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true
	// assign value
	bv, err := driver.Bool.ConvertValue(value)
	if err == nil {
		n.Bool = bv.(bool)
	}
	return err
}

// Value implements the driver Valuer interface.
func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil) // TODO: marshal to false?
}

func (n *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		n.Valid = true
		n.Bool = *b
	} else {
		n.Valid = false
	}
	return nil
}
