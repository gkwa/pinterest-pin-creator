package pinterest

import (
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse(`"2006-01-02T15:04:05"`, s)
	if err != nil {
		return err
	}
	*ct = CustomTime{t}
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format(`"2006-01-02T15:04:05"`)), nil
}
