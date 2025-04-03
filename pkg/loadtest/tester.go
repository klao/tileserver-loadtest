package loadtest

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Tester performs the load testing of a tile server
type Tester struct {
	config        Config
	tileGenerator *TileGenerator
	client        *http.Client
}

// NewTester creates a new Tester with the given configuration
func NewTester(config Config) *Tester {
	// Set up HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second, // Hard-coded 30 second timeout
	}

	return &Tester{
		config:        config,
		tileGenerator: NewTileGenerator(config),
		client:        client,
	}
}

// Run executes the load test according to the configuration
func (t *Tester) Run() error {
	// Validate the configuration
	if err := t.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	fmt.Println("Starting load test...")
	fmt.Printf("URL Template: %s\n", t.config.URLTemplate)
	fmt.Printf("Zoom: %d-%d, X: %d-%d, Y: %d-%d\n",
		t.config.MinZoom, t.config.MaxZoom,
		t.config.MinX, t.config.MaxX,
		t.config.MinY, t.config.MaxY)
	fmt.Printf("Threads: %d, Pattern: %s\n", t.config.Threads, t.config.Pattern)
	fmt.Printf("Name: %s, Environment: %s\n", t.config.Name, t.config.Environment)
	fmt.Printf("Output: %s\n", t.config.OutputPath)

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Set up context with cancellation for duration limit
	var ctx context.Context
	var cancel context.CancelFunc

	if t.config.Duration != "" {
		duration, err := time.ParseDuration(t.config.Duration)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
		ctx, cancel = context.WithTimeout(context.Background(), duration)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	// Channel for collecting results
	resultChan := make(chan Result, t.config.Threads*10)

	// Start worker goroutines
	fmt.Printf("Starting %d worker threads...\n", t.config.Threads)
	for i := 0; i < t.config.Threads; i++ {
		go t.worker(ctx, resultChan)
	}

	// Start metrics collection
	metrics := NewMetrics()
	metrics.Start()

	// Process results in the main thread
	requestCount := 0
	metrics.Start()

	// Process results until context is cancelled
	for {
		select {
		case <-ctx.Done():
			// Test duration exceeded
			fmt.Println("Test duration reached, stopping test...")
			metrics.End()

			// Don't close the channel here since workers might still be running
			// Just break out of the loop and let deferred cancel() terminate workers
			return t.writeResults(metrics.Results())

		case result, ok := <-resultChan:
			if !ok {
				// Channel closed, all workers done
				metrics.End()
				return t.writeResults(metrics.Results())
			}

			metrics.AddResult(result)
			requestCount++

			// Periodically print progress
			if requestCount%100 == 0 {
				fmt.Printf("Processed %d requests...\n", requestCount)
			}
		}
	}
}

// worker is a goroutine that makes requests to the tile server
func (t *Tester) worker(ctx context.Context, resultChan chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Generate tile and make request
			tile := t.tileGenerator.NextTile()
			url := t.tileGenerator.FormatURL(tile)

			startTime := time.Now()
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				resultChan <- Result{
					Latency:   0,
					Success:   false,
					Timestamp: time.Now(),
				}
				continue
			}

			resp, err := t.client.Do(req)
			latency := time.Since(startTime)

			if err != nil || resp.StatusCode != http.StatusOK {
				resultChan <- Result{
					Latency:   latency,
					Success:   false,
					Timestamp: time.Now(),
				}
			} else {
				resp.Body.Close()
				resultChan <- Result{
					Latency:   latency,
					Success:   true,
					Timestamp: time.Now(),
				}
			}
		}
	}
}

// validateConfig checks if the configuration is valid
func (t *Tester) validateConfig() error {
	if t.config.URLTemplate == "" {
		return errors.New("URL template is required")
	}

	if t.config.MinZoom > t.config.MaxZoom {
		return errors.New("min zoom must be less than or equal to max zoom")
	}

	if t.config.MinX > t.config.MaxX {
		return errors.New("min X must be less than or equal to max X")
	}

	if t.config.MinY > t.config.MaxY {
		return errors.New("min Y must be less than or equal to max Y")
	}

	if t.config.Pattern != "random" && t.config.Pattern != "fixed" {
		return errors.New("pattern must be either 'random' or 'fixed'")
	}

	if t.config.Duration != "" {
		_, err := time.ParseDuration(t.config.Duration)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
	}

	return nil
}

// writeResults writes the test results to the output CSV file
func (t *Tester) writeResults(results TestResults) error {
	// Check if file exists to determine if we need to write headers
	writeHeaders := false
	if _, err := os.Stat(t.config.OutputPath); os.IsNotExist(err) {
		writeHeaders = true
	}

	// Open file in append mode
	file, err := os.OpenFile(t.config.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Write headers if needed
	if writeHeaders {
		headers := "timestamp,name,environment,pattern,threads,min_zoom,max_zoom,min_x,max_x,min_y,max_y,total_requests,avg_latency,p95_latency,p99_latency,had_failures,duration\n"
		if _, err := file.WriteString(headers); err != nil {
			return fmt.Errorf("failed to write headers: %w", err)
		}
	}

	// Format current time
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Build CSV line
	fields := []string{
		timestamp,
		t.config.Name,
		t.config.Environment,
		t.config.Pattern,
		strconv.Itoa(t.config.Threads),
		strconv.Itoa(t.config.MinZoom),
		strconv.Itoa(t.config.MaxZoom),
		strconv.Itoa(t.config.MinX),
		strconv.Itoa(t.config.MaxX),
		strconv.Itoa(t.config.MinY),
		strconv.Itoa(t.config.MaxY),
		strconv.Itoa(results.TotalRequests),
		fmt.Sprintf("%.1f", results.AvgLatency),
		fmt.Sprintf("%.1f", results.P95Latency),
		fmt.Sprintf("%.1f", results.P99Latency),
		strconv.FormatBool(results.HadFailures),
		fmt.Sprintf("%.1f", results.TestDuration),
	}
	line := strings.Join(fields, ",") + "\n"

	// Write to file
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to write results: %w", err)
	}

	fmt.Println("Test completed successfully")
	fmt.Printf("Results written to %s\n", t.config.OutputPath)
	fmt.Printf("Total requests: %d\n", results.TotalRequests)
	fmt.Printf("Average latency: %.2f ms\n", results.AvgLatency)
	fmt.Printf("95th percentile: %.2f ms\n", results.P95Latency)
	fmt.Printf("99th percentile: %.2f ms\n", results.P99Latency)
	fmt.Printf("Success rate: %.2f%%\n", results.SuccessRate)
	fmt.Printf("Test duration: %.2f seconds\n", results.TestDuration)

	return nil
}
