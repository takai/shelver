package mover

import (
	"fmt"
	"os"
	"path/filepath"
)

type Move struct {
	Source string
	Dest   string
}

type Result struct {
	Moved   int
	Skipped int
	Failed  int
	Planned int
}

func ExecuteMoves(moves []Move, dryRun bool, destRoot string) Result {
	result := Result{
		Planned: len(moves),
	}
	
	for _, move := range moves {
		if dryRun {
			fmt.Printf("[DRYRUN] %s -> %s\n", move.Source, filepath.Join(destRoot, move.Dest))
			continue
		}
		
		status, err := executeMove(move, destRoot)
		if err != nil {
			fmt.Printf("ERROR: Failed to move %s: %v\n", move.Source, err)
			result.Failed++
		} else if status == "skipped" {
			result.Skipped++
		} else {
			result.Moved++
		}
	}
	
	return result
}

func executeMove(move Move, destRoot string) (string, error) {
	srcPath := move.Source
	destPath := filepath.Join(destRoot, move.Dest)
	
	// Check if source and destination are the same file
	if sameFile, err := isSameFile(srcPath, destPath); err == nil && sameFile {
		fmt.Printf("SKIP: %s (source and destination are the same)\n", srcPath)
		return "skipped", nil
	}
	
	// Create destination directory if needed
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}
	
	// Move file (this will overwrite if destination exists)
	if err := os.Rename(srcPath, destPath); err != nil {
		return "", fmt.Errorf("failed to rename %s to %s: %w", srcPath, destPath, err)
	}
	
	return "moved", nil
}

func isSameFile(src, dest string) (bool, error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return false, err
	}
	
	destInfo, err := os.Stat(dest)
	if err != nil {
		return false, err
	}
	
	return os.SameFile(srcInfo, destInfo), nil
}