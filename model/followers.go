package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Follower struct {
	FullName string `json:"full_name"`
	UserName string `json:"username"`
}

type Followers []Follower

// Value Marshal
func (a Followers) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Followers) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
