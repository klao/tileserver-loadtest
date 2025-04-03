# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands
- Build: `go build -o tile-load-test cmd/main.go`
- Run tests: `go test ./...`
- Run single test: `go test ./pkg/... -run TestName`
- Lint: `golangci-lint run`
- Format code: `gofmt -w .`

## Code Style Guidelines
- Use standard Go formatting conventions with `gofmt`
- IMPORTANT: Go code is indented with tabs, not spaces
- Follow Go naming conventions (CamelCase for exported, camelCase for internal)
- Imports: group standard library, external, and internal imports
- Use meaningful variable/function names that reflect purpose
- Add context to errors: `fmt.Errorf("doing X: %w", err)`
- Prefer strong typing over interface{} where possible
- Functions should have a single responsibility
- Use comments for exported functions and complex logic
- Utilize goroutines responsibly with proper synchronization
- Always check errors and handle them appropriately