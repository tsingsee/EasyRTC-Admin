package db

import (
	"time"

	"github.com/gocraft/dbr/v2"
)

const timeFormat = "2006-01-02 15:04:05.000000"

var Local = time.UTC

// NullTime is a type that can be null or a time.
type NullTime struct {
	dbr.NullTime
}

// NewNullTime creates a NullTime with Scan().
func NewNullTime(v interface{}) (n NullTime) {
	n.Scan(v)
	return
}

func (n *NullTime) Scan(value interface{}) error {
	var err error

	if value == nil {
		n.Time, n.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		n.Time, n.Valid = v, true
		return nil
	case []byte:
		n.Time, err = parseDateTime(string(v), Local)
		n.Valid = (err == nil)
		return err
	case string:
		n.Time, err = parseDateTime(v, Local)
		n.Valid = (err == nil)
		return err
	}

	n.Valid = false
	return nil
}

func parseDateTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26:
		if str == base[:len(str)] {
			return t, err
		}
		t, err = time.ParseInLocation(timeFormat[:len(str)], str, loc)
	default:
		err = dbr.ErrInvalidTimestring
	}

	return
}
