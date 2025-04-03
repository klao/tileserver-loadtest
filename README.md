# Tileserver Load Testing Tool

A command-line tool for load testing map tile servers, written in Go.

## Overview

This tool helps evaluate performance characteristics of tile servers by generating configurable request patterns and measuring response metrics.

## Usage

```
tile-load-test --url "https://tile.server/path/{z}/{x}/{y}.pbf"
               --min-zoom 10 --max-zoom 14
               --min-x 1000 --max-x 1100
               --min-y 1000 --max-y 1100
               --threads 10
               --pattern random
               --duration 60s
               --name "tileserver-ng"
               --environment "nginx+ram"
               --output results.csv
```

## Features

- Generate random or fixed pattern tile requests
- Configurable coordinate ranges and zoom levels
- Multi-threaded request handling
- CSV output for performance metrics
- Timeout and error handling

## Build

```
go build -o tile-load-test cmd/main.go
```

## License

This project is available under the MIT License.