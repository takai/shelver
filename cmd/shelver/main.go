package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"shelver/internal/mover"
	"shelver/internal/parser"
)

func main() {
	var (
		prefix = flag.String("prefix", "", "Group files by text before the specified prefix")
		dest   = flag.String("dest", ".", "Destination root directory")
		dryRun = flag.Bool("dryrun", false, "Show planned moves without executing")
	)
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: shelver [OPTIONS] <file/glob>...\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var files []string
	for _, arg := range flag.Args() {
		matches, err := filepath.Glob(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error expanding glob %s: %v\n", arg, err)
			continue
		}
		files = append(files, matches...)
	}

	if len(files) == 0 {
		fmt.Println("No files found matching the provided patterns")
		return
	}

	var moves []mover.Move
	var skipped int

	for _, file := range files {
		// Skip directories
		if info, err := os.Stat(file); err != nil || info.IsDir() {
			continue
		}

		// Use just the filename for parsing, not the full path
		filename := filepath.Base(file)
		group, ok := parser.ParseGroup(filename, *prefix)
		if !ok {
			skipped++
			continue
		}

		destPath := filepath.Join(group, filename)
		moves = append(moves, mover.Move{
			Source: file,
			Dest:   destPath,
		})
	}

	if len(moves) == 0 {
		fmt.Printf("No files matched the grouping patterns. Skipped: %d\n", skipped)
		return
	}

	result := mover.ExecuteMoves(moves, *dryRun, *dest)

	// Print summary
	if *dryRun {
		fmt.Printf("\nDry run complete. Would move %d files.\n", result.Planned)
	} else {
		fmt.Printf("\nMoved: %d, Skipped: %d, Failed: %d\n", result.Moved, result.Skipped, result.Failed)
	}

	if skipped > 0 {
		fmt.Printf("Files that didn't match pattern: %d\n", skipped)
	}
}
