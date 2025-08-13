# shelver

**A command-line file organizer that groups files into directories based on filename patterns.**

`shelver` automatically detects numeric sequences and optional prefixes in filenames to intelligently organize your files into logical groups. Perfect for organizing photo collections, music libraries, document series, and any other files with systematic naming patterns.

## Features

- **Automatic grouping** by trailing numeric sequences (1+ digits)
- **Prefix-based grouping** with user-specified prefixes
- **Dry-run mode** to preview changes before execution
- **Smart conflict handling** with overwrite behavior (like `mv`)
- **Unicode support** for international filenames
- **Safe operation** with same-file detection

## Installation

### Prerequisites

- Go 1.24.4 or later

### Build from Source

```bash
git clone <repository-url>
cd shelver

# Using Makefile (recommended)
make build

# Or using Go directly
go build -o shelver ./cmd/shelver
```

### Install

```bash
# Using Makefile - Install to system (requires sudo)
make install

# Using Makefile - Install to user directory (no sudo)
make install-user

# Or using Go directly
go install ./cmd/shelver
```

## Quick Start

```bash
# Group files by numeric sequences
shelver *.wav
# album-001.wav → album/album-001.wav
# album-002.wav → album/album-002.wav

# Group files by prefix
shelver --prefix p *.jpg
# kyoto-trip p04.jpg → kyoto-trip/kyoto-trip p04.jpg
# kyoto-trip p03.jpg → kyoto-trip/kyoto-trip p03.jpg

# Preview changes without moving files
shelver --dryrun *.txt
```

## Command-Line Interface

### Syntax

```bash
shelver [OPTIONS] <file/glob>...
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `--prefix <string>` | Group files by text before the specified prefix | *none* |
| `--dest <dir>` | Destination root directory | `.` |
| `--dryrun` | Show planned moves without executing | `false` |

### Examples

**Basic numeric grouping:**
```bash
shelver *.wav
# Input:  album-1.wav, album-2.wav, album-001.wav
# Output: album/album-1.wav, album/album-2.wav, album/album-001.wav
```

**Prefix-based grouping:**
```bash
shelver --prefix p *.jpg
# Input:  vacation p1.jpg, vacation p2.jpg
# Output: vacation/vacation p1.jpg, vacation/vacation p2.jpg
```

**Custom destination:**
```bash
shelver --dest /media/organized *.mp3
# Moves files to /media/organized/group/filename
```

**Dry run (preview only):**
```bash
shelver --dryrun --prefix page *.pdf
# Shows what would happen without moving files
```

## How It Works

### Organizing Photo Collections

```bash
# Before
ls
# kyoto-trip p1.jpg  kyoto-trip p2.jpg  kyoto-trip p3.jpg
# tokyo p1.jpg       tokyo p2.jpg

shelver --prefix p *.jpg

# After
ls
# kyoto-trip/
#   kyoto-trip p1.jpg
#   kyoto-trip p2.jpg
#   kyoto-trip p3.jpg
# tokyo/
#   tokyo p1.jpg
#   tokyo p2.jpg
```

### Organizing Music Files

```bash
# Before
ls
# album-1.wav  album-2.wav  album-3.wav
# mixtape-1.wav mixtape-2.wav

shelver *.wav

# After
ls
# album/
#   album-1.wav
#   album-2.wav
#   album-3.wav
# mixtape/
#   mixtape-1.wav
#   mixtape-2.wav
```

### Preview Mode

```bash
shelver --dryrun *.pdf
# [DRYRUN] report-1.pdf -> report/report-1.pdf
# [DRYRUN] report-2.pdf -> report/report-2.pdf
# [DRYRUN] manual-v2.pdf -> manual/manual-v2.pdf
#
# Dry run complete. Would move 3 files.
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/parser
go test ./internal/mover
```

## License

MIT