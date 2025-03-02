package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type VehicleStop struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Vehicle struct {
	Id              string      `json:"id"`
	Label           string      `json:"label"`
	Bearing         int         `json:"bearing"`
	DirectionId     int         `json:"directionId"`
	DirectionName   string      `json:"directionName"`
	TripStatus      string      `json:"tripStatus"`
	Latitude        float64     `json:"lat"`
	Longitude       float64     `json:"lng"`
	OccupancyStatus int         `json:"occupancyStatus"`
	TripId          string      `json:"tripId"`
	ShapeId         string      `json:"shapeId"`
	RouteId         string      `json:"routeId"`
	ServiceDate     string      `json:"serviceDate"`
	Headsign        string      `json:"headsign"`
	Timestamp       int         `json:"timestamp"`
	PreviousStop    VehicleStop `json:"prevStop"`
	CurrentStop     VehicleStop `json:"currentStop"`
	NextStop        VehicleStop `json:"nextStop"`
}

func (c *NextRideClient) GetRouteVehicles(routeId string) ([]Vehicle, error) {
	path := fmt.Sprintf("nextride/routes/%s/vehicles", routeId)
	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vehicles []Vehicle
	if err := json.NewDecoder(resp.Body).Decode(&vehicles); err != nil {
		slog.Error("Failed to decode response body", slog.Any("err", err))
		return nil, err
	}
	return vehicles, nil
}
