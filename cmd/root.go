package cmd

import (
	"sort"
)

var CLI struct {
	Route    RouteCmd    `cmd:"cmd" help:"Get route information"`
	Vehicles VehiclesCmd `cmd:"cmd" help:"Get vehicle information"`
}

type CliContext struct {
}

func sortRows(headers []string, rows [][]string, sortBy string) [][]string {
	// Find the index of the column to sort by
	sortIndex := -1
	for i, header := range headers {
		if header == sortBy {
			sortIndex = i
			break
		}
	}

	if sortIndex == -1 {
		// If the sortBy column is not found, return the rows as is
		return rows
	}

	// Sort the rows based on the sortIndex column
	sort.Slice(rows, func(i, j int) bool {
		if rows[i][sortIndex] == rows[j][sortIndex] {
			return rows[i][0] < rows[j][0]
		}
		return rows[i][sortIndex] < rows[j][sortIndex]
	})

	return rows
}
