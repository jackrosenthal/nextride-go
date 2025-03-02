package api

import (
	"encoding/json"
	"log/slog"
)

type RouteAlertInfo struct {
	Name                  string `json:"routeName"`
	LatestAlertExpiration int    `json:"latestAlertExpiration"`
	AlertCount            int    `json:"alertCount"`
	RouteId               string `json:"routeId"`
	NewestAlertCreatedAt  int    `json:"newestAlertCreatedAt"`
	DisplayName           string `json:"displayName"`
}

type StationAlertInfo struct {
	Id           int    `json:"id"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	Header       string `json:"header"`
	Cause        string `json:"cause"`
	DateCreated  int    `json:"dateCreated"`
	DateModified int    `json:"dateModified"`
	DelayTime    string `json:"delayTime"`
	Severity     int    `json:"severity"`
	StartDate    int    `json:"startDate"`
	Link         *struct {
		Text string `json:"text"`
		To   string `json:"to"`
	} `json:"link"`
	EndDate        *int    `json:"endDate"`
	AlertLifecycle string  `json:"alertLifecycle"`
	Relevance      float64 `json:"relevance"`
}

type AlertListsAttributes struct {
	Bus     []RouteAlertInfo   `json:"bus"`
	Rail    []RouteAlertInfo   `json:"rail"`
	Station []StationAlertInfo `json:"station"`
}

type alertListsData struct {
	Attributes AlertListsAttributes `json:"attributes"`
}

type alertListsResponse struct {
	Data alertListsData `json:"data"`
}

func (c *NextRideClient) GetAlerts() (*AlertListsAttributes, error) {
	resp, err := c.Get("rider-alerts/alerts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var alerts alertListsResponse
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		slog.Error("Failed to decode response body", slog.Any("err", err))
		return nil, err
	}
	return &alerts.Data.Attributes, nil
}
