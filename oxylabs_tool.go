package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OxylabsToolSet wraps Oxylabs Scraper API configurations.
type OxylabsToolSet struct {
	username string
	password string
}

// NewOxylabsToolSet initializes a new OxylabsToolSet.
func NewOxylabsToolSet(username, password string) *OxylabsToolSet {
	return &OxylabsToolSet{
		username: username,
		password: password,
	}
}

// ScrapeWebPage extracts real-time content from any public webpage using Oxylabs Universal Scraper.
func (o *OxylabsToolSet) ScrapeWebPage(ctx context.Context, url string) (string, error) {
	reqPayload := map[string]any{
		"source": "universal",
		"url":    url,
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal scrape request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://realtime.oxylabs.io/v1/queries", bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(o.username, o.password)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("oxylabs scrape failed (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return string(bodyBytes), nil
}

// SearchWeb queries google search using Oxylabs Search Scraper.
func (o *OxylabsToolSet) SearchWeb(ctx context.Context, query string) (string, error) {
	reqPayload := map[string]any{
		"source": "google_search",
		"query":  query,
		"context": []map[string]any{
			{
				"key":   "filter",
				"value": 1,
			},
		},
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search request: %w", err)
	}

	// Using the endpoint provided in the curl sample: https://data.oxylabs.io/v1/queries
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://data.oxylabs.io/v1/queries", bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(o.username, o.password)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("oxylabs search failed (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return string(bodyBytes), nil
}
