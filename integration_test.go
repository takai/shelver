package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestShelver_Integration(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"album-001.wav",
		"album-002.wav",
		"kyoto-trip p04.jpg",
		"kyoto-trip p03.jpg",
		"sp02.jpg",
	}

	for _, filename := range testFiles {
		fullPath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Build shelver binary for testing
	cmd := exec.Command("go", "build", "-o", filepath.Join(tempDir, "shelver"), "./cmd/shelver")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build shelver: %v", err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	t.Run("numeric grouping", func(t *testing.T) {
		// Test numeric grouping
		cmd := exec.Command("./shelver", "*.wav", "sp02.jpg")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("shelver command failed: %v\nOutput: %s", err, output)
		}

		// Verify files were moved to correct groups
		albumDir := filepath.Join(tempDir, "album")
		spDir := filepath.Join(tempDir, "sp")

		expectedFiles := []string{
			filepath.Join(albumDir, "album-001.wav"),
			filepath.Join(albumDir, "album-002.wav"),
			filepath.Join(spDir, "sp02.jpg"),
		}

		for _, expectedFile := range expectedFiles {
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
				t.Errorf("Expected file not found: %s", expectedFile)
			}
		}

		// Verify source files were moved (no longer exist in original location)
		for _, filename := range []string{"album-001.wav", "album-002.wav", "sp02.jpg"} {
			if _, err := os.Stat(filename); !os.IsNotExist(err) {
				t.Errorf("Source file still exists: %s", filename)
			}
		}
	})

	t.Run("prefix grouping", func(t *testing.T) {
		// Test prefix grouping on remaining files
		cmd := exec.Command("./shelver", "--prefix", "p", "*.jpg")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("shelver command failed: %v\nOutput: %s", err, output)
		}

		// Verify files were moved to kyoto-trip group
		kyotoDir := filepath.Join(tempDir, "kyoto-trip")
		expectedFiles := []string{
			filepath.Join(kyotoDir, "kyoto-trip p03.jpg"),
			filepath.Join(kyotoDir, "kyoto-trip p04.jpg"),
		}

		for _, expectedFile := range expectedFiles {
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
				t.Errorf("Expected file not found: %s", expectedFile)
			}
		}
	})
}

func TestShelver_DryRun(t *testing.T) {
	tempDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tempDir, "test-99.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Build shelver binary
	cmd := exec.Command("go", "build", "-o", filepath.Join(tempDir, "shelver"), "./cmd/shelver")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build shelver: %v", err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Run dry run
	cmd = exec.Command("./shelver", "--dryrun", "test-99.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("shelver dryrun failed: %v\nOutput: %s", err, output)
	}

	// Verify dry run output
	if !strings.Contains(string(output), "[DRYRUN]") {
		t.Errorf("Expected [DRYRUN] in output, got: %s", output)
	}

	if !strings.Contains(string(output), "test/test-99.txt") {
		t.Errorf("Expected move plan in output, got: %s", output)
	}

	// Verify file was not actually moved
	if _, err := os.Stat("test-99.txt"); os.IsNotExist(err) {
		t.Errorf("File was moved during dry run")
	}

	if _, err := os.Stat("test"); !os.IsNotExist(err) {
		t.Errorf("Directory was created during dry run")
	}
}
