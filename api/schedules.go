package api

import (
	"encoding/json"
	"log/slog"
	"slices"
)

type Runboard struct {
	Id                   string `json:"id"`
	CurrentId            string `json:"currentId"`
	CurrentName          string `json:"currentName"`
	CurrentEffectiveDate string `json:"currentEffectiveDate"`
}

type ModeType int

const (
	BusMode ModeType = iota
	RailMode
	FlexRideMode
)

type RouteType int

const (
	DenverLocalBus     RouteType = 1
	LightRail          RouteType = 2
	FreeMallRideBus    RouteType = 3
	DenverRegionalBus  RouteType = 5
	AirportBus         RouteType = 9
	BoulderLocalBus    RouteType = 11
	BoulderRegionalBus RouteType = 12
	LongmontLocalBus   RouteType = 15
	FlatironFlyerBus   RouteType = 30
	RegionalRail       RouteType = 32
)

func (r RouteType) String() string {
	switch r {
	case DenverLocalBus:
		return "Denver Local Bus"
	case LightRail:
		return "Light Rail"
	case FreeMallRideBus:
		return "Free Mall Ride Bus"
	case DenverRegionalBus:
		return "Denver Regional Bus"
	case AirportBus:
		return "Airport Bus"
	case BoulderLocalBus:
		return "Boulder Local Bus"
	case BoulderRegionalBus:
		return "Boulder Regional Bus"
	case LongmontLocalBus:
		return "Longmont Local Bus"
	case FlatironFlyerBus:
		return "Flatiron Flyer Bus"
	case RegionalRail:
		return "Regional Rail"
	default:
		return "Unknown"
	}
}

func (r RouteType) GetModeType() ModeType {
	switch r {
	case DenverLocalBus, DenverRegionalBus, AirportBus, BoulderLocalBus, BoulderRegionalBus, LongmontLocalBus, FlatironFlyerBus:
		return BusMode
	case LightRail, RegionalRail:
		return RailMode
	default:
		return 0
	}
}

type RouteRegion int

const (
	DenverRegion   RouteRegion = 1
	BoulderRegion  RouteRegion = 2
	LongmontRegion RouteRegion = 3
	AirportRegion  RouteRegion = 4
)

func (r RouteType) AsRegion() RouteRegion {
	switch r {
	case DenverLocalBus, DenverRegionalBus, FreeMallRideBus:
		return DenverRegion
	case AirportBus:
		return AirportRegion
	case BoulderLocalBus, BoulderRegionalBus, FlatironFlyerBus:
		return BoulderRegion
	case LongmontLocalBus:
		return LongmontRegion
	default:
		return 0
	}
}

type Route struct {
	Id        string    `json:"id"`
	Branch    *string   `json:"branch"`
	RouteId   string    `json:"routeId"`
	RouteType RouteType `json:"routeType"`
	RouteName string    `json:"routeName"`
	LineName  string    `json:"lineName"`
}

func (r *Route) IsInRegion(region RouteRegion) bool {
	if r.Id == "FF3" && region == BoulderRegion {
		return false
	}

	if r.RouteType.AsRegion() == region {
		return true
	}

	switch region {
	case DenverRegion:
		if r.RouteType == FlatironFlyerBus {
			return true
		}
		if r.Id == "AT/ATA" {
			return true
		}
	case LongmontRegion:
		if r.Id == "BOLT" {
			return true
		}
	case BoulderRegion:
		if r.Id == "AB1/AB2" {
			return true
		}
	case AirportRegion:
		return slices.Contains([]string{"104L", "145X", "169L", "A"}, r.Id)
	}

	return false
}

type RunboardResponse struct {
	Id       string   `json:"id"`
	Runboard Runboard `json:"runboard"`
	Routes   []Route  `json:"routes"`
}

func (c *NextRideClient) GetCurrentRunboard() (*RunboardResponse, error) {
	resp, err := c.Get("schedules")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var runboard RunboardResponse
	if err := json.NewDecoder(resp.Body).Decode(&runboard); err != nil {
		slog.Error("Failed to decode response body", slog.Any("err", err))
		return nil, err
	}
	return &runboard, nil
}
