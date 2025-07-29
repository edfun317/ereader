package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProgressStore stores progress for multiple books
type ProgressStore struct {
	Progresses map[string]ReadingProgress `json:"progresses"`
}

// ReadingProgress stores the reading position for a book
type ReadingProgress struct {
	FilePath string     `json:"file_path"`
	Position CurrentPos `json:"position"`
}

// saveProgress saves the current reading position to a file
func (v *CLIViewer) saveProgress() error {
	if v.currentFile == "" {
		return fmt.Errorf("no current file set")
	}

	// Create progress directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	progressDir := filepath.Join(homeDir, ".ereader")
	if err := os.MkdirAll(progressDir, 0755); err != nil {
		return fmt.Errorf("failed to create progress directory: %w", err)
	}

	progressFile := filepath.Join(progressDir, "progress.json")

	// Read existing progress store or create new one
	store := ProgressStore{
		Progresses: make(map[string]ReadingProgress),
	}

	// Try to read existing progress file
	if data, err := os.ReadFile(progressFile); err == nil {
		if err := json.Unmarshal(data, &store); err != nil {
			// If the file is corrupted, start fresh
			store.Progresses = make(map[string]ReadingProgress)
		}
	}

	// Update or add new progress
	store.Progresses[v.currentFile] = ReadingProgress{
		FilePath: v.currentFile,
		Position: v.currentPos,
	}

	// Marshal the entire store
	data, err := json.MarshalIndent(store, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal progress data: %w", err)
	}

	// Write to temporary file first
	tempFile := progressFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary progress file: %w", err)
	}

	// Rename temporary file to actual file (atomic operation)
	if err := os.Rename(tempFile, progressFile); err != nil {
		os.Remove(tempFile) // Clean up temp file if rename fails
		return fmt.Errorf("failed to save progress file: %w", err)
	}

	return nil
}

// loadProgress loads the reading position from file
func (v *CLIViewer) loadProgress() error {
	if v.currentFile == "" {
		return fmt.Errorf("no current file set")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	progressFile := filepath.Join(homeDir, ".ereader", "progress.json")
	data, err := os.ReadFile(progressFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No saved progress, start from beginning
		}
		return fmt.Errorf("failed to read progress file: %w", err)
	}

	var store ProgressStore
	if err := json.Unmarshal(data, &store); err != nil {
		return fmt.Errorf("failed to unmarshal progress data: %w", err)
	}

	// Load progress for current file if it exists
	if progress, exists := store.Progresses[v.currentFile]; exists {
		// Validate the position before setting
		if progress.Position.Chapter >= 0 &&
			progress.Position.Chapter < v.reader.GetTotalChapters() {
			v.currentPos = progress.Position
		}
	}

	return nil
}

// Debug function to help troubleshoot progress saving
func (v *CLIViewer) debugProgress() {
	homeDir, _ := os.UserHomeDir()
	progressFile := filepath.Join(homeDir, ".ereader", "progress.json")

	fmt.Printf("\nDebug Progress Information:\n")
	fmt.Printf("Current File: %s\n", v.currentFile)
	fmt.Printf("Current Position: Chapter %d, Page %d\n",
		v.currentPos.Chapter, v.currentPos.Page)
	fmt.Printf("Progress File: %s\n", progressFile)

	if data, err := os.ReadFile(progressFile); err == nil {
		fmt.Printf("Current Progress File Content:\n%s\n", string(data))
	} else {
		fmt.Printf("Error reading progress file: %v\n", err)
	}
}
