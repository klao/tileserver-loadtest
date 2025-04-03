# Map Tile Server Load Testing Tool Specification

## Overview
A command-line tool for load testing map tile servers, written in Go. This tool helps evaluate performance characteristics of tile servers by generating configurable request patterns and measuring response metrics.

## Core Requirements

### Input Parameters
- **Tile URL Template**: A URL pattern with placeholders for zoom, x, and y coordinates (e.g., `https://tile.server/path/{z}/{x}/{y}.pbf`)
- **Coordinate Ranges**:
  - Zoom level (`z`)
  - X coordinates: Min and max values
  - Y coordinates: Min and max values
- **Concurrency**: Number of parallel threads to use for requests
- **Name**: Identifier for the tile server being tested (e.g., "tileserver-ng")
- **Environment**: Description of the server environment (e.g., "nginx+ram")

### Performance Metrics
- **Average Latency**: Mean response time across all requests
- **High Percentile Latency**: 95th and 99th percentile response times
- **Failures**: Report if any request failures occurred during testing

## Additional Requirements

### Request Patterns
The tool should support two request patterns:
1. **Random**: Randomly select tiles within the specified ranges
2. **Fixed**: Repeatedly request the same set of tiles (for hammering specific tiles)

### Test Duration
- **Duration Flag**: Allow setting a maximum test duration (e.g., run for 60 seconds)
- The test should stop after either the duration expires or all requested tiles have been processed

### Error Handling
- **Timeout**: Use a hard-coded timeout value (e.g., 30 seconds) for requests
- **Failure Reporting**: Report if any failures occurred but without detailed statistics

### Output Format
- **CSV Output**: Each test should append a single line to a CSV file
- The CSV should include all test parameters and results
- Format should include: timestamp, name, environment, request pattern, thread count, zoom, min/max x/y, total requests, average latency, 95th percentile, 99th percentile, failures flag, test duration

## Implementation Details

### Language and Dependencies
- Implement in Go
- Use standard library where possible
- Minimal external dependencies

### Command-Line Interface
```
tile-load-test --url 'https://tile.server/path/{z}/{x}/{y}.pbf'
               --zoom 14
               --min-x 1000 --max-x 1100
               --min-y 1000 --max-y 1100
               --pattern random
               --threads 10
               --duration 60s
               --name 'tileserver-ng'
               --environment 'nginx+ram'
               --output results.csv
```

### CSV Format
Example CSV line:
```
2025-04-03T14:30:45Z,tileserver-ng,nginx+ram,random,10,14,1000,1100,1000,1100,5000,127.3,245.6,389.2,false,60.0
```
Fields:
1. Timestamp
2. Name
3. Environment
4. Pattern
5. Thread count
6. Zoom
7. Min X
8. Max X
9. Min Y
10. Max Y
11. Total requests
12. Average latency (ms)
13. 95th percentile latency (ms)
14. 99th percentile latency (ms)
15. Had failures (true/false)
16. Test duration (seconds)

### Error Handling
- Validate input parameters before starting the test
- Gracefully handle network errors and server failures
- Support clean termination via interrupt signals (Ctrl+C)
- Create CSV file if it doesn't exist, with appropriate headers

### Performance Considerations
- Efficiently manage goroutines to avoid resource exhaustion
- Minimize memory usage when tracking metrics
- Ensure accurate timing measurements

## Future Extensions
These are not part of the initial requirements but worth considering for future versions:
- Request rate limiting
- Authentication support
- Detailed failure analysis
- Multiple output formats beyond CSV
- Resumable test configurations
- More sophisticated request patterns
