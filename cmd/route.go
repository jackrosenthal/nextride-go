package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jamespfennell/gtfs"
)

type RouteListCmd struct {
	SortBy string `option:"" name:"sort" help:"Sort by field" default:"Type"`
}

func (c *RouteListCmd) Run(context *CliContext) error {
	resp, err := http.Get("https://www.rtd-denver.com/files/gtfs/google_transit.zip")
	if err != nil {
		slog.Error("Failed to get GTFS data", slog.Any("err", err))
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", slog.Any("err", err))
		return err
	}

	gtfsData, err := gtfs.ParseStatic(respBody, gtfs.ParseStaticOptions{})
	if err != nil {
		slog.Error("Failed to parse GTFS data", slog.Any("err", err))
		return err
	}

	rows := [][]string{}
	for _, route := range gtfsData.Routes {
		rows = append(rows, []string{
			route.Id,
			route.ShortName,
			route.LongName,
			route.Description,
			route.Type.String(),
		})
	}

	headers := []string{"ID", "Short Name", "Long Name", "Description", "Type"}
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
