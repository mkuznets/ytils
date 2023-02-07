package ytime

import (
	"database/sql/driver"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	t.Time.IsZero()
	return t.Time.MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	return t.Time.UnmarshalJSON(data)
}

func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.UnixMilli(), nil
}

func (t *Time) Scan(src interface{}) error {
	if src == nil {
		t.Time = time.Time{}
		return nil
	}
	t.Time = time.UnixMilli(src.(int64))
	return nil
}

func New(t time.Time) Time {
	return Time{t}
}

func Now() Time {
	return New(time.Now())
}
