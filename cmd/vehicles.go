package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/jackrosenthal/nextride-go/api"
)

type VehiclesCmd struct {
	SortBy string   `option:"" name:"sort" help:"Sort by field" default:"Vehicle Number"`
	Routes []string `arg:"" name:"route" help:"Route IDs to get vehicle information for"`
}

func (c *VehiclesCmd) Run(context *CliContext) error {
	allVehicles := []api.Vehicle{}
	for _, route := range c.Routes {
		vehicles, err := context.Client.GetRouteVehicles(route)
		if err != nil {
			return err
		}
		allVehicles = append(allVehicles, vehicles...)
	}

	rows := make([][]string, len(allVehicles))
	for i, vehicle := range allVehicles {
		rows[i] = []string{
			vehicle.Label,
			vehicle.RouteId,
			vehicle.DirectionName,
			vehicle.Headsign,
			fmt.Sprintf("%d", vehicle.OccupancyStatus),
			fmt.Sprintf("%f, %f", vehicle.Latitude, vehicle.Longitude),
			vehicle.CurrentStop.Name,
			vehicle.NextStop.Name,
		}
	}

	headers := []string{"Vehicle Number", "Route", "Direction", "Headsign", "Occupancy", "Location", "Current Stop", "Next Stop"}

	rows = sortRows(headers, rows, c.SortBy)

	tab := table.New().
		Border(lipgloss.NormalBorder()).
		Headers(headers...).
		Rows(rows...)

	fmt.Println(tab)
	return nil
}
