package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/klao/tileserver-loadtest/pkg/loadtest"
)

func main() {
	// Parse command line flags
	url := flag.String("url", "", "Tile URL template (e.g., https://tile.server/path/{z}/{x}/{y}.pbf)")
	zoom := flag.Int("zoom", 0, "Zoom level")
	minX := flag.Int("min-x", 0, "Minimum X coordinate")
	maxX := flag.Int("max-x", 0, "Maximum X coordinate")
	minY := flag.Int("min-y", 0, "Minimum Y coordinate")
	maxY := flag.Int("max-y", 0, "Maximum Y coordinate")
	threads := flag.Int("threads", 1, "Number of concurrent threads")
	pattern := flag.String("pattern", "random", "Request pattern: random or fixed")
	duration := flag.String("duration", "", "Maximum test duration (e.g., 60s)")
	name := flag.String("name", "default", "Name identifier for the tile server")
	environment := flag.String("environment", "default", "Environment description")
	output := flag.String("output", "results.csv", "Output CSV file path")

	flag.Parse()

	// Validate required parameters
	if *url == "" {
		fmt.Println("Error: --url parameter is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create and run the test
	config := loadtest.Config{
		URLTemplate: *url,
		Zoom:        *zoom,
		MinX:        *minX,
		MaxX:        *maxX,
		MinY:        *minY,
		MaxY:        *maxY,
		Threads:     *threads,
		Pattern:     *pattern,
		Duration:    *duration,
		Name:        *name,
		Environment: *environment,
		OutputPath:  *output,
	}

	tester := loadtest.NewTester(config)
	if err := tester.Run(); err != nil {
		fmt.Printf("Error running test: %v\n", err)
		os.Exit(1)
	}
}
