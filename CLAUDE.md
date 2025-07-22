# Titleparser Development Guide

**Important**: This project follows standardized guidelines from the `llm-shared/` submodule:

- **General project standards**: `llm-shared/project_tech_stack.md`
- **Go-specific guidelines**: `llm-shared/languages/go.md`
- **GitHub issue management**: `llm-shared/GITHUB.md`
- **Function analysis tools**: `llm-shared/utils/` (gofuncs, validate-docs)

## Architecture Overview

Titleparser is a Go AWS Lambda function that extracts titles from URLs with custom parsers for specific sites (Reddit, YouTube, HackerNews, etc.). The architecture uses a plugin-style pattern where handlers register themselves via `init()` functions.

**Core Flow:**

1. `lambda/main.go` receives TitleQuery via AWS Lambda
2. Iterates through registered handlers (pattern matching against URL)
3. Falls back to default OpenGraph/HTML title extraction
4. Caches results in DynamoDB (disabled in local mode)

## Build Commands

- Build (Lambda): `task build` - Creates ARM64 Linux binary + ZIP for AWS
- Build (Local): `task build-local` - Native binary for testing
- Test All: `task test` - Runs all tests with coverage
- Run Single Test: `go test -v ./handler -run TestSpecificFunction`
- Lint: `task lint` - golangci-lint (continues on errors)
- Clean: `task clean`

## Local Testing

The application supports local testing mode via the `RUNMODE` environment variable:

**1. Build and run locally:**

```bash
task build-local          # Creates native binary 'titleparser'
RUNMODE=local ./titleparser  # Starts HTTP server on localhost:8081
```

**2. Test with stdin mode (fastest for development):**

```bash
# Simple URL test:
echo '{"url":"https://example.com"}' | RUNMODE=stdin ./titleparser

# Test custom handler (Reddit):
echo '{"url":"https://reddit.com/r/golang/comments/123/test"}' | RUNMODE=stdin ./titleparser

# From file:
cat test-query.json | RUNMODE=stdin ./titleparser
```

**3. Test the HTTP server endpoint:**

```bash
# Using httpie:
http POST http://localhost:8081/title url="http://example.com"

# Using curl:
curl -X POST http://localhost:8081/title \
  -H "Content-Type: application/json" \
  -d '{"url": "http://reddit.com/r/golang/comments/123/test"}'
```

**4. Available endpoints (HTTP server mode):**

- `POST /title` - Main title parsing (same JSON format as Lambda)
- `GET /hi` - Simple health check
- `GET /` - Hello world test

**Local mode differences:**

- **RUNMODE=stdin**: Reads JSON from stdin, outputs to stdout, then exits (fastest for testing)
- **RUNMODE=local**: Starts HTTP server on localhost:8081 for interactive testing
- Both modes skip DynamoDB caching (results not stored)
- All custom handlers work identically to production

## Handler Registration Pattern

Each site parser follows this pattern (see `handler/reddit.go`):

```go
var RedditMatch = `.*reddit\.com/r/.*/comments/.*/.*`  // Package-level regex

func Reddit(url string) (string, error) {
    // Custom parsing logic
}

func init() {
    lambda.RegisterHandler(RedditMatch, Reddit)  // Auto-registration
}
```

## Code Style Guidelines

**Critical formatting requirements (from llm-shared/languages/go.md):**

- **Always run `goimports -w .` after changes** (NOT gofmt - goimports includes gofmt + import management)
- **Always run `task build` before finishing** - includes tests, linting, and compilation
- Install prerequisites: `go install golang.org/x/tools/cmd/goimports@latest`

**Project-specific patterns:**

- Use table-driven tests with `t.Parallel()` for concurrency
- Define regex constants at package level (e.g., `RedditMatch`)
- Structured logging with logrus (`log.Infof`, `log.Warnf`)
- HTTP requests use common headers from `handler/config.go`
- Prefer OpenGraph titles (`og:title`) over HTML `<title>`
- Error handling: return descriptive errors, don't use `log.Fatal` in handlers
- Test files use `_test.go` suffix with testdata fixtures
- **Task completion requires**: passing tests, successful linting, and successful build

## Key Implementation Details

- **Local vs Lambda mode**: Check `RUNMODE` env var to skip DynamoDB caching
- **HTTP client**: 10-second timeout, custom User-Agent from `handler/config.go`
- **Title sanitization**: Max 200 chars, whitespace cleanup in `lambda/default.go`
- **Handler priority**: First regex match wins, default handler is fallback
- **Error types**: Use specific errors like `ErrTitleNotFound`, `ErrNotHTML`
- **Function analysis**: Use `go run llm-shared/utils/gofuncs.go -dir .` to list all functions
- **Project validation**: Use `go run llm-shared/utils/validate-docs.go` to check project structure

## Project Structure

- `handler/`: Site-specific parsers (reddit.go, youtube.go, etc.)
- `lambda/`: Core functionality (main.go, default.go, cache.go)
- `build/`: Generated artifacts (bootstrap binary, ZIP file)
- `llm-shared/`: Git submodule with shared standards
- `.cursor/rules/`: Project-specific coding rules
