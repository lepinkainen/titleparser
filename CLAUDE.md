# Titleparser Development Guide

## Build Commands
- Build (Lambda): `task build`
- Build (Local): `task build-local`
- Test All: `task test`
- Run Single Test: `go test -v ./handler -run TestSpecificFunction`
- Lint: `task lint`
- Clean: `task clean`

## Code Style Guidelines
- Follow Go standard practices (gofmt formatting)
- Use table-driven tests with t.Parallel() for concurrency
- Define regex constants at package level
- Proper error handling with specific error types
- Structured logging
- Register handlers via init() functions
- Use descriptive function and variable names
- Prefer OpenGraph titles when available
- Use testdata fixtures for consistent testing
- Check `.cursor/rules/` directory for project-specific Cursor rules

## Project Structure
- `handler/`: Site-specific parsers (Reddit, YouTube, HackerNews, etc.)
- `lambda/`: Core functionality (request handling, caching)
- `.cursor/rules/`: Cursor rule files with project standards and guidelines