# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

**Build the project:**
```bash
make build          # Optimized build to build/shelver
make build-dev      # Development build (faster)
make build-local    # Build to ./shelver
```

**Testing:**
```bash
make test           # Run all tests
make test-verbose   # Run tests with verbose output
make test-coverage  # Run tests with coverage report
go test ./...       # Direct go test command
go test -v ./internal/parser    # Test specific package
go test -v ./internal/mover     # Test specific package
```

**Code quality:**
```bash
make check          # Run fmt, vet, and test
make fmt            # Format code
make vet            # Vet code
go fmt ./...        # Direct go fmt command
go vet ./...        # Direct go vet command
```

**Running the application:**
```bash
make run-example    # Run with testdata examples in dry-run mode
./build/shelver --dryrun testdata/*.wav testdata/*.jpg
```

## Project Architecture

`shelver` is a Go CLI tool with a clean separation of concerns:

- **`cmd/shelver/main.go`**: CLI entry point, flag parsing, orchestrates parser and mover
- **`internal/parser/`**: Core filename parsing logic for extracting group names
- **`internal/mover/`**: File moving operations and conflict handling

### Key Architecture Patterns

**Parser Package (`internal/parser/`):**
- `ParseGroup(filename, prefix)` is the main entry point
- Two matching forms: prefix-form (when `--prefix` is set) and numeric-only form
- Uses regex patterns: `numericPattern` for 2+ digit sequences, custom patterns for prefix matching
- Boundary protection prevents false matches (e.g., "sp02" with `--prefix p`)

**Mover Package (`internal/mover/`):**
- `ExecuteMoves()` handles both dry-run and actual execution
- `Move` struct represents source → destination file operations
- Same-file detection prevents moving files to themselves
- Directory creation handled automatically

**Main Application Flow:**
1. Parse CLI flags (`--prefix`, `--dest`, `--dryrun`)
2. Expand glob patterns to file list
3. For each file: parse group name using `parser.ParseGroup()`
4. Build list of `Move` operations
5. Execute moves via `mover.ExecuteMoves()`

### Important Implementation Details

**Filename Processing:**
- Only the base filename (not full path) is used for parsing
- Filename stem (without extension) is extracted for pattern matching
- Group names have trailing separators trimmed

**Regex Patterns:**
- Numeric pattern: `^(.+?)[-_.\s]*(\d+)$` (requires 1+ digits)
- Prefix boundary checking when no separator present
- Escaped prefix strings to handle special regex characters

**Testing Structure:**
- Unit tests in each package (`*_test.go`, `edge_cases_test.go`)
- Integration tests at root level (`integration_test.go`)
- Uses standard Go testing with `t.TempDir()` for file operations
- Test data in `testdata/` directory

## Development Guidelines

- Follow TDD principles: write tests first, then implement
- Use `make check` before committing to ensure code quality
- Integration tests should use temporary directories, not modify testdata
- All file operations should be tested in both dry-run and actual modes
- Error handling should be consistent with Unix `mv` behavior
- Do not add footer on commit message
