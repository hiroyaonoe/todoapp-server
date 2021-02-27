package entity

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

/*
NullString はsql.NullStringをうまくJSONにMarshal/Unmarshal出来るようにするstruct
(https://okamuuu.hatenablog.com/entry/2016/12/20/150339)
名前付きフィールドにすることでStringerインターフェースに対応させる．
(https://stackoverflow.com/questions/65559059/custom-golang-sql-nullstring-stringer-interface)
*/
type NullString struct {
	ns sql.NullString
}

func (s *NullString) Scan(value interface{}) error {
	return s.ns.Scan(value)
}

func (s NullString) Value() (driver.Value, error) {
	return s.ns.Value()
}

func NewNullString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: s != ""}}
}

func (s *NullString) Set(str string) {
	s.ns.String = str
	s.ns.Valid = str != ""
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.ns.Valid {
		return json.Marshal(s.String())
	} else {
		return json.Marshal(nil)
	}
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	s.Set(str)
	return nil
}

func (s NullString) String() string {
	return s.ns.String
}

func (s *NullString) IsNull() bool {
	return !s.ns.Valid
}

func (s NullString) Equal(t NullString) bool {
	if s.ns.Valid {
		return t.ns.Valid
	} else {
		return s.String() == t.String()
	}
}
