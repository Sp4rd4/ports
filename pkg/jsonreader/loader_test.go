package jsonreader_test

import (
	"io"
	"strings"
	"testing"

	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/jsonreader"
	"github.com/stretchr/testify/assert"
)

var bufferSize = 512

var examplesLoad = []struct {
	name   string
	reader io.Reader
	result []*domain.Port
}{
	{
		name:   "Empty reader",
		reader: strings.NewReader(""),
		result: nil,
	},
	{
		name:   "Nil reader",
		result: nil,
	},
	{
		name:   "Incorrect json",
		reader: strings.NewReader(`{"name":"Pretoria","coordinates":[28.22,-25.7],"city":"Pretoria","province":"Gauteng","country":"South Africa","alias":[],"regions":[],"timezone":"Africa/Johannesburg","unlocs":["ZAPRY"]}`),
		result: nil,
	},
	{
		name:   "Correct json",
		reader: strings.NewReader(`{"AEAJM":{"name":"Ajman","city":"Ajman","country":"United Arab Emirates","alias":[],"regions":[],"coordinates":[55.5136433,25.4052165],"province":"Ajman","timezone":"Asia/Dubai","unlocs":["AEAJM"],"code":"52000"},"ZAPLZ":{"name":"Port Elizabeth","city":"Port Elizabeth","country":"South Africa","alias":[],"regions":[],"coordinates":[25.5207358,-33.7139247],"province":"Eastern Cape","timezone":"Africa/Johannesburg","unlocs":["ZAPLZ"],"code":"79145"}}`),
		result: []*domain.Port{
			&domain.Port{
				ID:      "AEAJM",
				Name:    "Ajman",
				City:    "Ajman",
				Country: "United Arab Emirates",
				Alias:   domain.StringArray{},
				Regions: domain.StringArray{},
				Coordinates: domain.Location{
					Latitude:  55.5136433,
					Longitude: 25.4052165,
				},
				Province: "Ajman",
				Timezone: "Asia/Dubai",
				Unlocs:   domain.StringArray{"AEAJM"},
				Code:     "52000",
			},
			&domain.Port{
				ID:      "ZAPLZ",
				Name:    "Port Elizabeth",
				City:    "Port Elizabeth",
				Country: "South Africa",
				Alias:   domain.StringArray{},
				Regions: domain.StringArray{},
				Coordinates: domain.Location{
					Latitude:  25.5207358,
					Longitude: -33.7139247,
				},
				Province: "Eastern Cape",
				Timezone: "Africa/Johannesburg",
				Unlocs:   domain.StringArray{"ZAPLZ"},
				Code:     "79145",
			},
		},
	},
}

func TestLoad(t *testing.T) {
	for _, ex := range examplesLoad {
		loader := jsonreader.NewLoader(ex.reader, bufferSize, nil)
		t.Run(ex.name, func(t *testing.T) {
			var res []*domain.Port
			c := loader.Load()
			for p := range c {
				res = append(res, p)
			}
			assert.ElementsMatch(t, res, ex.result, "Chanel should return expected ports")
		})
	}
}
