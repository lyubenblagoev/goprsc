package goprsc

import "time"

// DateTime represents a time that can be unmarshaled from a JSON string
// formatted as "yyyy-MM-dd HH:mm:ss". All exported methods of time.Time
// can be called on DateTime.
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
	t.Time, err = time.Parse(`"2006-01-02 15:04:05"`, str)
	return err
}

// Equal reports whether t and u are equal based on time.Equal().
func (t DateTime) Equal(u DateTime) bool {
	return t.Time.Equal(u.Time)
}
