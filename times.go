package codecov

import (
	"time"
)

const (
	// RFC3339UTC represents RFC3339 but with an assumed UTC timestamp
	RFC3339UTC = "2006-01-02T15:04:05"

	// RFC3339UTCSansT represents RFC3339UTC but without the "T" separating character
	RFC3339UTCSansT = "2006-01-02 15:04:05"

	// RFC3339NanoSansT represents RFC3339Nano but without the "T" separating character
	RFC3339NanoSansT = "2006-01-02 15:04:05.999999Z07:00"
)

// Time is a time object that handles custom formats coming from Codecov
// When marshaling and unmarshaling it will make multiple attempts to parse in formats
// that have been observed from the Codecov API.
//
// Examples:
// * "2006-01-02T15:04:05"
// * "2006-01-02 15:04:05"
type Time struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in one of the formats supported by either
// the standard UnmarshalJSON call or one of the constant formats defined in this package.
func (t *Time) UnmarshalJSON(data []byte) error {
	strData := string(data)

	// Ignore null, like in the main JSON package.
	if strData == "null" {
		return nil
	}

	var err error
	var tm time.Time
	for _, format := range []string{time.RFC3339, RFC3339UTC, RFC3339UTCSansT, RFC3339NanoSansT} {
		if tm, err = time.Parse(`"`+format+`"`, strData); err == nil {
			t.Time = tm.UTC()
			return err
		}
	}

	return err
}
