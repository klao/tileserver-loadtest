package loadtest

// Config holds the configuration parameters for a load test
type Config struct {
	URLTemplate  string // URL pattern with placeholders for zoom, x, and y coordinates
	MinZoom      int    // Minimum zoom level
	MaxZoom      int    // Maximum zoom level
	MinX         int    // Minimum X coordinate
	MaxX         int    // Maximum X coordinate
	MinY         int    // Minimum Y coordinate
	MaxY         int    // Maximum Y coordinate
	Threads      int    // Number of concurrent threads
	Pattern      string // Request pattern: "random" or "fixed"
	Duration     string // Maximum test duration (e.g., "60s")
	Name         string // Identifier for the tile server
	Environment  string // Description of the server environment
	OutputPath   string // Path to the output CSV file
}
