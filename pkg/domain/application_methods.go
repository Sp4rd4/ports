package domain

import (
	"database/sql/driver"
	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
)

var (
	//nolint
	json = jsoniter.ConfigDefault

	ErrTypeAssertion = errors.New("type assertion .([]byte) failed")
	ErrUnmarshal     = errors.New("byte unmarshal failed")
)

// MarshalJSON implements custom marshal.
func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]float64{l.Latitude, l.Longitude})
}

// UnmarshalJSON implements custom unmarshal.
func (l *Location) UnmarshalJSON(b []byte) error {
	ll := [2]float64{}
	err := json.Unmarshal(b, &ll)
	if err != nil {
		return err
	}
	l.Latitude, l.Longitude = ll[0], ll[1]
	return nil
}

// Scan implements the sql.Scanner interface.
func (l *Location) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return ErrTypeAssertion
	}

	err := json.Unmarshal(source, l)
	if err != nil {
		return ErrUnmarshal
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (l Location) Value() (driver.Value, error) {
	j, err := json.Marshal(l)
	return j, err
}

// Scan implements the sql.Scanner interface.
func (a *StringArray) Scan(src interface{}) error {
	return (*pq.StringArray)(a).Scan(src)
}

// Value implements the driver.Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}
