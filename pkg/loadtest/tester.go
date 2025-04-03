package loadtest

import (
	"errors"
	"fmt"
	"time"
)

// Tester performs the load testing of a tile server
type Tester struct {
	config Config
}

// NewTester creates a new Tester with the given configuration
func NewTester(config Config) *Tester {
	return &Tester{
		config: config,
	}
}

// Run executes the load test according to the configuration
func (t *Tester) Run() error {
	// Validate the configuration
	if err := t.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// TODO: Implement the actual load testing logic
	fmt.Println("Starting load test...")
	fmt.Printf("URL Template: %s\n", t.config.URLTemplate)
	fmt.Printf("Zoom: %d-%d, X: %d-%d, Y: %d-%d\n",
		t.config.MinZoom, t.config.MaxZoom,
		t.config.MinX, t.config.MaxX,
		t.config.MinY, t.config.MaxY)
	fmt.Printf("Threads: %d, Pattern: %s\n", t.config.Threads, t.config.Pattern)
	fmt.Printf("Name: %s, Environment: %s\n", t.config.Name, t.config.Environment)
	fmt.Printf("Output: %s\n", t.config.OutputPath)

	// Placeholder for the actual test implementation
	time.Sleep(1 * time.Second)

	return nil
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

	return nil
}
