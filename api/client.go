package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type NextRideClient struct {
	client *http.Client
}

func NewNextRideClient() *NextRideClient {
	return &NextRideClient{
		client: &http.Client{},
	}
}

type NextRideError struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

type NextRideErrorsResponse struct {
	Errors []NextRideError `json:"errors"`
}

func (c *NextRideClient) Get(path string) (*http.Response, error) {
	url := fmt.Sprintf("https://nodejs-prod.rtd-denver.com/api/v2/%s", path)
	resp, err := c.client.Get(url)
	if err != nil {
		slog.Error("Failed to fetch data from NextRide API", slog.Any("url", url), slog.Any("err", err))
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()

		var errors NextRideErrorsResponse
		if err := json.NewDecoder(resp.Body).Decode(&errors); err != nil {
			slog.Error("Failed to decode error response body", slog.Any("err", err))
			return nil, err
		}

		slog.Error("Failed to fetch data from NextRide API", slog.Any("url", url), slog.Any("errors", errors))
		return nil, fmt.Errorf("failed to fetch data from NextRide API: %s", errors.Errors[0].Detail)
	}

	return resp, nil
}
