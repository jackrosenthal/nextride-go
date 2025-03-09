package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/jamespfennell/gtfs"
	"github.com/jamespfennell/gtfs/proto"
)

type VehiclesCmd struct {
	SortBy string   `option:"" name:"sort" help:"Sort by field" default:"Vehicle Number"`
	Routes []string `arg:"" name:"route" help:"Route IDs to get vehicle information for" optional:""`
}

func (c *VehiclesCmd) Run(context *CliContext) error {
	resp, err := http.Get("https://www.rtd-denver.com/files/gtfs-rt/VehiclePosition.pb")
	if err != nil {
		slog.Error("Failed to get vehicle positions", slog.Any("err", err))
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", slog.Any("err", err))
		return err
	}
	rtData, err := gtfs.ParseRealtime(respBody, &gtfs.ParseRealtimeOptions{})
	if err != nil {
		slog.Error("Failed to parse realtime data", slog.Any("err", err))
		return err
	}

	allVehicles := []gtfs.Vehicle{}
	for _, vehicle := range rtData.Vehicles {
		trip := vehicle.GetTrip()
		if len(c.Routes) == 0 || slices.Contains(c.Routes, trip.ID.RouteID) {
			allVehicles = append(allVehicles, vehicle)
		}
	}

	rows := make([][]string, len(allVehicles))
	for i, vehicle := range allVehicles {
		trip := vehicle.GetTrip()
		pos := vehicle.Position
		occupancy := vehicle.OccupancyStatus
		if occupancy == nil {
			occ := proto.VehiclePosition_NO_DATA_AVAILABLE
			occupancy = &occ
		}

		odoStr := "???"
		if pos.Odometer != nil {
			odoStr = fmt.Sprintf("%f mi", *pos.Odometer/1609.34)
		}

		speedStr := "???"
		if pos.Speed != nil {
			speedStr = fmt.Sprintf("%f MPH", *pos.Speed*2.23694)
		}

		stopStr := "None"
		if vehicle.StopID != nil {
			stopStr = *vehicle.StopID
		}

		rows[i] = []string{
			vehicle.GetID().Label,
			trip.ID.RouteID,
			trip.ID.DirectionID.String(),
			trip.ID.ID,
			occupancy.String(),
			fmt.Sprintf("%f, %f", *pos.Latitude, *pos.Longitude),
			odoStr,
			speedStr,
			stopStr,
		}
	}

	headers := []string{"Vehicle Number", "Route", "Direction", "Trip ID", "Occupancy", "Location", "Odometer", "Speed", "Stop ID"}

	rows = sortRows(headers, rows, c.SortBy)

	tab := table.New().
		Border(lipgloss.NormalBorder()).
		Headers(headers...).
		Rows(rows...)

	fmt.Println(tab)
	return nil
}
