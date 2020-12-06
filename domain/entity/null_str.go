package entity

import (
	"database/sql"
	"encoding/json"
)

/*
NullString はsql.NullStringをうまくJSONにMarshal/Unmarshal出来るようにするstruct
(https://okamuuu.hatenablog.com/entry/2016/12/20/150339)
*/
type NullString struct {
	sql.NullString
}

func NewNullString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: s != ""}}
}

func (s *NullString) Set(str string) {
	s.String = str
	s.Valid = str != ""
	// new := NewNullString(str)
	// s = &new
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	} else {
		return json.Marshal(nil)
	}
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	s.String = str
	s.Valid = str != ""
	return nil
}

func (s *NullString) ToString() string {
	if s.IsNull() {
		return ""
	}
	return s.String
}

func (s *NullString) IsNull() bool {
	return !s.Valid
}

func (s NullString) Equals(t NullString) bool {
	if s.Valid {
		return t.Valid
	} else {
		return s.String == t.String
	}
}
