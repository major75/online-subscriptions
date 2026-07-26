package types

import (
	"errors"
	"time"
)

const DATE_FORMAT = "01-2006"

type MMYYYYDate struct {
	time.Time
}

func (d *MMYYYYDate) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == `""` {
		return errors.New("empty date value")
	}

	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	t, err := time.Parse(DATE_FORMAT, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func (d *MMYYYYDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(DATE_FORMAT) + `"`), nil
}
