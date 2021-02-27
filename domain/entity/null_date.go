package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	layout = "2006-01-02"
)

type NullDate struct {
	sql.NullTime
}

func NewNullDate(s string) NullDate {
	date, err := newNullDate(s)
	if err != nil {
		panic(err.Error())
	}
	return *date
}

func (d *NullDate) Set(str string) error {
	new, err := newNullDate(str)
	if err != nil {
		return err
	}
	d.Time = new.Time
	d.Valid = new.Valid
	return nil
}

func newNullDate(s string) (*NullDate, error) {
	array := strings.Split(s, "-")

	if len(array) != 3 {
		return nil, fmt.Errorf("%s is invalid format", s)
	}

	ymd := []int{-1, -1, -1}
	var err error
	for i, v := range array {
		ymd[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	y, m, d := ymd[0], ymd[1], ymd[2]

	date, err := isExist(y, m, d)
	if err != nil {
		return nil, err
	}

	res := new(NullDate)
	res.Time = date
	res.Valid = !date.IsZero()

	return res, nil
}

func (d *NullDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *NullDate) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	err := d.Set(str)
	return err
}

func (d NullDate) GetTime() time.Time {
	if d.IsNull() {
		return time.Unix(0, 0)
	}
	return d.Time
}

func (d NullDate) String() string {
	return d.GetTime().Format(layout)
}

func (d NullDate) IsNull() bool {
	return !d.Valid
}

func isExist(year, month, day int) (time.Time, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if date.Year() == year && date.Month() == time.Month(month) && date.Day() == day {
		return date, nil
	} else {
		return time.Time{}, fmt.Errorf("%d-%d-%d is not exist", year, month, day)
	}
}

func (s NullDate) Equal(t NullDate) bool {
	return s.String() == t.String() && s.IsNull() == t.IsNull()
}
