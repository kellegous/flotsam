package util

import "time"

type Time struct {
	Year      int    `json:"year" jsonschema:"year"`
	Month     int    `json:"month" jsonschema:"month"`
	Day       int    `json:"day" jsonschema:"day"`
	Hour      int    `json:"hour" jsonschema:"hour"`
	Minute    int    `json:"minute" jsonschema:"minute"`
	Second    int    `json:"second" jsonschema:"second"`
	DayOfWeek string `json:"day_of_week" jsonschema:"day_of_week"`
	Timezone  string `json:"timezone" jsonschema:"timezone"`
}

func ToTime(t time.Time, tz string) (*Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}

	t = t.In(loc)

	return &Time{
		Year:      t.Year(),
		Month:     int(t.Month()),
		Day:       t.Day(),
		Hour:      t.Hour(),
		Minute:    t.Minute(),
		Second:    t.Second(),
		DayOfWeek: t.Weekday().String(),
		Timezone:  loc.String(),
	}, nil
}
