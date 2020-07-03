package domain

import (
	jsoniter "github.com/json-iterator/go"
)

type Port struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	City        string   `json:"city" db:"city"`
	Country     string   `json:"country" db:"country"`
	Alias       []string `json:"alias" db:"alias"`
	Regions     []string `json:"regions" db:"regions"`
	Coordinates Location `json:"coordinates" db:"coordinates"`
	Province    string   `json:"province" db:"province"`
	Timezone    string   `json:"timezone" db:"timezone"`
	Unlocs      []string `json:"unlocs" db:"unlocs"`
	Code        string   `json:"code" db:"code"`
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type PortRepository interface {
	Save(port *Port) error
	Get(id string) (*Port, error)
}

var json = jsoniter.ConfigFastest

func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]float64{l.Latitude, l.Longitude})
}

func (l *Location) UnmarshalJSON(b []byte) error {
	ll := [2]float64{}
	err := json.Unmarshal(b, &ll)
	if err != nil {
		return err
	}

	l.Latitude, l.Longitude = ll[0], ll[1]
	return nil
}
