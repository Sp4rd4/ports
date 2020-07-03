package domain

type StringArray []string

type Port struct {
	ID          string      `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	City        string      `json:"city" db:"city"`
	Country     string      `json:"country" db:"country"`
	Alias       StringArray `json:"alias" db:"alias"`
	Regions     StringArray `json:"regions" db:"regions"`
	Coordinates Location    `json:"coordinates" db:"coordinates"`
	Province    string      `json:"province" db:"province"`
	Timezone    string      `json:"timezone" db:"timezone"`
	Unlocs      StringArray `json:"unlocs" db:"unlocs"`
	Code        string      `json:"code" db:"code"`
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type PortRepository interface {
	Save(port *Port) error
	Get(id string) (*Port, error)
}
