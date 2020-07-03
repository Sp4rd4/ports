package proto

import "github.com/sp4rd4/ports/pkg/domain"

func PortDomainToProto(p *domain.Port) *Port {
	if p == nil {
		return &Port{}
	}
	return &Port{
		Id:   p.ID,
		Name: p.Name,
		City: p.City,
		Code: p.Code,
		Coordinates: &Location{
			Latitude:  p.Coordinates.Latitude,
			Longitude: p.Coordinates.Longitude,
		},
		Country:  p.Country,
		Alias:    p.Alias,
		Regions:  p.Regions,
		Province: p.Province,
		Timezone: p.Timezone,
		Unlocs:   p.Unlocs,
	}
}

func PortProtoToDomain(p *Port) *domain.Port {
	if p == nil {
		return &domain.Port{}
	}
	port := &domain.Port{
		ID:       p.Id,
		Name:     p.Name,
		City:     p.City,
		Code:     p.Code,
		Country:  p.Country,
		Alias:    p.Alias,
		Regions:  p.Regions,
		Province: p.Province,
		Timezone: p.Timezone,
		Unlocs:   p.Unlocs,
	}
	if p.Coordinates != nil {
		port.Coordinates = domain.Location{
			Latitude:  p.Coordinates.Latitude,
			Longitude: p.Coordinates.Longitude,
		}
	}
	return port
}
