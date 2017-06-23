package goprsc

import "time"

// DateTime represents a time that can be unmarshaled from a JSON string
// formatted as "yyyy-mm-ddThh:mm:ss+|-hhmm" (e.g. '2017-01-02T15:47:59+0100').
// All exported methods of time.Time can be called on DateTime.
type DateTime struct {
	time.Time
}

func (t DateTime) String() string {
	return t.Time.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *DateTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	var err error
	t.Time, err = time.Parse(`"2006-01-02T15:04:05-0700"`, str)
	return err
}

// Equal reports whether t and u are equal based on time.Equal().
func (t DateTime) Equal(u DateTime) bool {
	return t.Time.Equal(u.Time)
}
