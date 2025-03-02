package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jackrosenthal/nextride-go/api"
)

type RouteListCmd struct {
	SortBy string `option:"" name:"sort" help:"Sort by field" default:"Type"`
	Mode   string `option:"" name:"mode" help:"Filter by mode (bus or rail)" default:"" enum:"bus,rail,"`
	Region string `option:"" name:"region" help:"Filter by region" default:"" enum:"denver,boulder,longmont,airport,"`
}

func routeRegionFromString(region string) api.RouteRegion {
	switch region {
	case "denver":
		return api.DenverRegion
	case "boulder":
		return api.BoulderRegion
	case "longmont":
		return api.LongmontRegion
	case "airport":
		return api.AirportRegion
	}
	return 0
}

func (c *RouteListCmd) Run(context *CliContext) error {
	runboard, err := context.Client.GetCurrentRunboard()
	if err != nil {
		return err
	}

	rows := [][]string{}
	for _, route := range runboard.Routes {
		if c.Mode != "" && c.Mode != strings.ToLower(route.RouteType.String()) {
			continue
		}
		if c.Region != "" && !route.IsInRegion(routeRegionFromString(c.Region)) {
			continue
		}
		rows = append(rows, []string{
			route.Id,
			route.RouteId,
			route.RouteName,
			route.RouteType.String(),
		})
	}

	headers := []string{"Name", "ID", "Description", "Type"}
	rows = sortRows(headers, rows, c.SortBy)
	table := table.New().
		Border(lipgloss.NormalBorder()).
		Headers(headers...).
		Rows(rows...)

	fmt.Println(table)

	return nil
}

type RouteCmd struct {
	List RouteListCmd `cmd:"" help:"List all routes"`
}
