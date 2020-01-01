package lunch_server

import (
	"encoding/json"
	"fmt"
)

type Date struct {
	y, m, d int
}

func NewDate(year, month, day int) Date {
	return Date{year, month, day}
}

func (d Date) Year() int {
	return d.y
}

func (d Date) Month() int {
	return d.m
}

func (d Date) Day() int {
	return d.d
}

func (d Date) Before(other Date) bool {
	if d.Year() < other.Year() {
		return true
	}
	if d.Year() > other.Year() {
		return false
	}
	if d.Month() < other.Month() {
		return true
	}
	if d.Month() > other.Month() {
		return false
	}
	return d.Day() < other.Day()
}

func (d Date) Equal(other Date) bool {
	if d.Year() != other.Year() {
		return false
	}
	if d.Month() != other.Month() {
		return false
	}
	return d.Day() == other.Day()
}

func (d Date) After(other Date) bool {
	return other.Before(d)
}

func (d Date) String() string {
	return fmt.Sprintf("%4d%2d%2d", d.y, d.m, d.d)
}

// MarshalJSON implements the Marshaller interface
func (d Date) MarshalJSON() ([]byte, error) {
	type D struct {
		Y int `json:"year"`
		M int `json:"month"`
		D int `json:"day"`
	}
	return json.Marshal(D{d.y, d.m, d.d})
}

// UnmarshalJSON implements the Unmarshaller interface
func (d *Date) UnmarshalJSON(bts []byte) error {
	type D struct {
		Y int `json:"year"`
		M int `json:"month"`
		D int `json:"day"`
	}
	var tmp D
	err := json.Unmarshal(bts, &tmp)
	if err != nil {
		return err
	}
	d.y, d.m, d.d = tmp.Y, tmp.M, tmp.D
	return nil
}