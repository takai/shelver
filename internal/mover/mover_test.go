package mover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMover_DryRun(t *testing.T) {
	moves := []Move{
		{
			Source: "album-001.wav",
			Dest:   "album/album-001.wav",
		},
		{
			Source: "album-002.wav", 
			Dest:   "album/album-002.wav",
		},
		{
			Source: "kyoto-trip p04.jpg",
			Dest:   "kyoto-trip/kyoto-trip p04.jpg",
		},
	}

	result := ExecuteMoves(moves, true, ".")
	
	if result.Moved != 0 {
		t.Errorf("DryRun should not move any files, got %d", result.Moved)
	}
	
	expectedPlanned := 3
	if result.Planned != expectedPlanned {
		t.Errorf("Expected %d planned moves, got %d", expectedPlanned, result.Planned)
	}
	
	if result.Failed != 0 {
		t.Errorf("DryRun should not have failures, got %d", result.Failed)
	}
}

func TestMover_ActualMoves(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files
	testFiles := []string{"album-001.wav", "album-002.wav"}
	for _, filename := range testFiles {
		fullPath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}
	
	moves := []Move{
		{
			Source: "album-001.wav",
			Dest:   "album/album-001.wav",
		},
		{
			Source: "album-002.wav",
			Dest:   "album/album-002.wav",
		},
	}
	
	// Change to temp directory for test
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)
	
	result := ExecuteMoves(moves, false, ".")
	
	if result.Moved != 2 {
		t.Errorf("Expected 2 moved files, got %d", result.Moved)
	}
	
	if result.Failed != 0 {
		t.Errorf("Expected 0 failed moves, got %d", result.Failed)
	}
	
	// Verify files were moved
	for _, move := range moves {
		destPath := filepath.Join(tempDir, move.Dest)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("File was not moved to %s", destPath)
		}
		
		srcPath := filepath.Join(tempDir, move.Source)
		if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
			t.Errorf("Source file still exists at %s", srcPath)
		}
	}
}

func TestMover_ConflictHandling(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create source file
	srcPath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(srcPath, []byte("source content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Create destination directory and file
	destDir := filepath.Join(tempDir, "group")
	os.MkdirAll(destDir, 0755)
	destPath := filepath.Join(destDir, "test.txt")
	if err := os.WriteFile(destPath, []byte("existing content"), 0644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	
	moves := []Move{
		{
			Source: "test.txt",
			Dest:   "group/test.txt",
		},
	}
	
	// Change to temp directory for test
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)
	
	result := ExecuteMoves(moves, false, ".")
	
	// Should overwrite existing file (like mv behavior)
	if result.Moved != 1 {
		t.Errorf("Expected 1 moved file, got %d", result.Moved)
	}
	
	// Verify content was overwritten
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	
	if string(content) != "source content" {
		t.Errorf("File was not overwritten. Expected 'source content', got '%s'", string(content))
	}
}

func TestMover_SameFileDetection(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test file
	filename := "test.txt"
	srcPath := filepath.Join(tempDir, filename)
	if err := os.WriteFile(srcPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	moves := []Move{
		{
			Source: filename,
			Dest:   filename, // Same file
		},
	}
	
	// Change to temp directory for test
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)
	
	result := ExecuteMoves(moves, false, ".")
	
	// Should skip same file moves
	if result.Moved != 0 {
		t.Errorf("Expected 0 moved files for same file, got %d", result.Moved)
	}
	
	if result.Skipped != 1 {
		t.Errorf("Expected 1 skipped file, got %d", result.Skipped)
	}
}