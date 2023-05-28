package monitor

import (
	"context"
	"net/http"
	"strings"
)

// PingResponse is the response from the Ping Endpoint
type PingResponse struct {
	Up bool `json:"up"`
}

// Ping pings a specific site and determines whether it's up or down right now.
//
//encore:api public path=/ping/*url
func Ping(ctx context.Context, url string) (*PingResponse, error) {
	// Check if the url starts with "http" or "https", default "https"
	if !strings.HasPrefix(url, "http:") && !strings.HasPrefix(url, "https:") {
		url = "https://" + url
	}

	// Make ping request to check if its up
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	// Status code lower than 400 is considered down
	up := resp.StatusCode < 400

	return &PingResponse{Up: up}, nil
}
