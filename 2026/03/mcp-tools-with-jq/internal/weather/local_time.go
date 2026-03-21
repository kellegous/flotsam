package weather

import (
	"encoding/json"
	"time"
)

type LocalTime time.Time

func (t LocalTime) ToTime() time.Time {
	return time.Time(t)
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	s := time.Time(t).Format(time.RFC3339)
	return json.Marshal(s)
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*t = LocalTime(tt.In(time.Local))
	return nil
}
