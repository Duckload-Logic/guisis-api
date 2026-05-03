package structs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type NullableTime sql.NullTime

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}

	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	nt.Time = t
	nt.Valid = true
	return nil
}

func (nt *NullableTime) Scan(value interface{}) error {
	var t sql.NullTime
	if err := t.Scan(value); err != nil {
		return err
	}
	nt.Time = t.Time
	nt.Valid = t.Valid
	return nil
}

func (nt NullableTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func FromSqlNullTime(nt sql.NullTime) NullableTime {
	return NullableTime{Time: nt.Time, Valid: nt.Valid}
}

func ToSqlNullTime(nt NullableTime) sql.NullTime {
	return sql.NullTime{Time: nt.Time, Valid: nt.Valid}
}

func TimeToNullableTime(t time.Time) NullableTime {
	if t.IsZero() {
		return NullableTime{Valid: false}
	}
	return NullableTime{Time: t, Valid: true}
}
