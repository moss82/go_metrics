package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// MetricData represents the parsed Prometheus metric data.
type MetricData struct {
	Name        string
	Description string
	Type        string
}

// Dashboard represents a basic Grafana dashboard.
type Dashboard struct {
	Title   string    `json:"title"`
	Panels  []Panel   `json:"panels"`
	UID     string    `json:"uid"`    // Unique identifier for the dashboard
	Version int       `json:"version"` // Version of the dashboard (needed for newer versions)
}

// Panel represents a Grafana panel.
type Panel struct {
	Title   string       `json:"title"`
	Type    string       `json:"type"`
	Targets []Target     `json:"targets"`
	GridPos GridPosition `json:"gridPos"` // Grid positioning for new Grafana panel layouts
}

// Target represents a query to be used in a panel.
type Target struct {
	Expr string `json:"expr"`
}

// GridPosition defines the panel's position and size in the dashboard.
type GridPosition struct {
	H int `json:"h"` // height
	W int `json:"w"` // width
	X int `json:"x"` // x position in the grid
	Y int `json:"y"` // y position in the grid
}

// FetchPrometheusMetrics fetches metrics from a remote Prometheus endpoint.
func FetchPrometheusMetrics(url string) ([]MetricData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	lines := strings.Split(string(body), "\n")
	metrics := []MetricData{}
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			// Ignore comments, types, and help descriptions.
			continue
		}
		// Example: parsing lines like `my_metric_name 123`
		parts := strings.Fields(line)
		if len(parts) > 0 {
			metrics = append(metrics, MetricData{Name: parts[0], Type: "gauge"})
		}
	}
	return metrics, nil
}

// GenerateGrafanaDashboard creates a simple Grafana dashboard from the parsed metrics.
func GenerateGrafanaDashboard(metrics []MetricData) Dashboard {
	panels := []Panel{}
	for i, metric := range metrics {
		panel := Panel{
			Title: metric.Name,
			Type:  "timeseries", // Use React-based "timeseries" panel instead of deprecated Angular graph panel
			Targets: []Target{
				{Expr: metric.Name},  // Set the PromQL expression to the metric name
			},
			GridPos: GridPosition{
				H: 8,  // Height of the panel
				W: 12, // Width of the panel (this determines layout in Grafana)
				X: (i % 2) * 12, // X position (alternate between 0 and 12)
				Y: (i / 2) * 8,  // Y position (stack vertically)
			},
		}
		panels = append(panels, panel)
	}
	dashboard := Dashboard{
		Title:   "Auto-Generated Dashboard",
		Panels:  panels,
		UID:     "auto-generated-dashboard", // Set a UID for the dashboard
		Version: 1,                          // Set dashboard version (required for new Grafana versions)
	}
	return dashboard
}

// SaveDashboardToFile saves the dashboard JSON to a local file.
func SaveDashboardToFile(dashboard Dashboard, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func main() {
	// Check if a Prometheus URL was provided as a command-line argument
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run main.go <Prometheus URL>")
	}
	prometheusURL := os.Args[1]

	// Fetch the Prometheus metrics
	metrics, err := FetchPrometheusMetrics(prometheusURL)
	if err != nil {
		log.Fatalf("Error fetching metrics: %v", err)
	}

	// Generate a Grafana dashboard
	dashboard := GenerateGrafanaDashboard(metrics)

	// Save the dashboard as a JSON file
	filename := "grafana_dashboard.json"
	err = SaveDashboardToFile(dashboard, filename)
	if err != nil {
		log.Fatalf("Error saving dashboard: %v", err)
	}

	fmt.Printf("Dashboard saved to %s\n", filename)
}
