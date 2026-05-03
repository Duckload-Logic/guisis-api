package structs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type (
	NullableString sql.NullString
	NullableInt64  sql.NullInt64
)

func (ns *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	ns.String = s
	ns.Valid = true
	return nil
}

func (ns *NullableString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	ns.String = s.String
	ns.Valid = s.Valid
	return nil
}

func (ns NullableString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns NullableString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func StringToNullableString(s string) NullableString {
	if s == "" {
		return NullableString{Valid: false}
	}
	return NullableString{String: s, Valid: true}
}

func PointerToNullableString(ps *string) NullableString {
	if ps == nil || *ps == "" {
		return NullableString{Valid: false}
	}
	return NullableString{String: *ps, Valid: true}
}

func FromSqlNull(ns sql.NullString) NullableString {
	return NullableString{String: ns.String, Valid: ns.Valid}
}

func ToSqlNull(ns NullableString) sql.NullString {
	return sql.NullString{String: ns.String, Valid: ns.Valid}
}

func (ni *NullableInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		ni.Int64 = 0
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	ni.Int64 = i
	ni.Valid = true
	return nil
}

func (ni *NullableInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}
	ni.Int64 = i.Int64
	ni.Valid = i.Valid
	return nil
}

func (ni NullableInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

func (ni NullableInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ni.Int64)
}

func FromSqlNullInt64(ni sql.NullInt64) NullableInt64 {
	return NullableInt64{Int64: ni.Int64, Valid: ni.Valid}
}

func ToSqlNullInt64(ni NullableInt64) sql.NullInt64 {
	return sql.NullInt64{Int64: ni.Int64, Valid: ni.Valid}
}
