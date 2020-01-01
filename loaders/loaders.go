package loaders

import (
	"fmt"
	"strings"
	"time"

	"github.com/snabb/isoweek"
)

// dayFromString converts Swedish days to their corresponding time.Weekday
func dayFromString(s string) (time.Weekday, error) {
	switch strings.ToLower(s) {
	case "måndag":
		return time.Monday, nil
	case "tisdag":
		return time.Tuesday, nil
	case "onsdag":
		return time.Wednesday, nil
	case "torsdag":
		return time.Thursday, nil
	case "fredag":
		return time.Friday, nil
	case "lördag":
		return time.Saturday, nil
	case "söndag":
		return time.Sunday, nil
	default:
		return -1, fmt.Errorf("%s is not a valid day", s)
	}
}

func parseISOWeek(year, week int, day time.Weekday) time.Time {
	t := isoweek.StartTime(year, week, time.UTC)
	t = t.AddDate(0, 0, (int(day)+6)%7)
	return t
}

func daysUntil(now, then time.Time) int {
	return int(then.Sub(now.Truncate(time.Hour*24)).Hours() / 24)
}
