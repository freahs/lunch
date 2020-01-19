package lunch

import (
	"encoding/json"
	"github.com/snabb/isoweek"
	"time"
)

type Date struct {
	t time.Time
}

func Now() Date {
	now := time.Now()
	return Date{time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)}
}

func Week(year, week int) Date {
	y, m, d := isoweek.StartDate(year, week)
	return Date{time.Date(y, m, d, 0, 0, 0, 0, time.Local)}
}

func NewDate(year, month, day int) Date {
	return Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)}
}

func (d Date) Year() int {
	return d.t.Year()
}

func (d Date) Month() int {
	return int(d.t.Month())
}

func (d Date) Day() int {
	return d.t.Day()
}

func (d Date) Week() (year, week int) {
	return d.t.ISOWeek()
}

func (d Date) Weekday() string {
	return d.t.Weekday().String()
}

func (d Date) Before(other Date) bool {
	return d.t.Before(other.t)
}

func (d Date) Equal(other Date) bool {
	return d.t.Equal(other.t)
}

func (d Date) After(other Date) bool {
	return other.Before(d)
}

func (d Date) Sub(other Date) time.Duration {
	return d.t.Sub(other.t)
}

func (d Date) Add(years, months, days int) Date {
	return Date{d.t.AddDate(years, months, days)}
}

func (d Date) String() string {
	return d.t.Format("Mon 2006-01-02")
}

// MarshalJSON implements the Marshaller interface
func (d Date) MarshalJSON() ([]byte, error) {
	s := d.t.Format("20060102")
	return json.Marshal(s)
}

// UnmarshalJSON implements the Unmarshaller interface
func (d *Date) UnmarshalJSON(bts []byte) error {
	var s string
	err := json.Unmarshal(bts, &s)
	if err != nil {
		return err
	}
	d.t, err = time.Parse("20060102", s)
	return err
}
