package entity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	layout = "2006-01-02"
)

type Date struct {
	date time.Time
}

func NewDate(s string) (*Date, error) {
	array := strings.Split(s, "-")

	if len(array) != 3 {
		return nil, fmt.Errorf("%s is invalid format", s)
	}

	var ymd []int
	var err error
	for i, v := range array {
		ymd[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	y, m, d := ymd[0], ymd[1], ymd[2]

	date, err := isExist(y, m, d)

	return &Date{date: date}, nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	d, err := NewDate(str)
	return err
}

func (d Date) String() string {
	return d.date.Format(layout)
}

func isExist(year, month, day int) (time.Time, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if date.Year() == year && date.Month() == time.Month(month) && date.Day() == day {
		return date, nil
	} else {
		return time.Time{}, fmt.Errorf("%d-%d-%d is not exist", year, month, day)
	}
}
