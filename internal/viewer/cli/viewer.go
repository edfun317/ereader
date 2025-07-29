package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

func (v *CLIViewer) eventLoop() error {
	for !v.shouldExit {
		if err := v.displayCurrentPage(); err != nil {
			return err
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			return fmt.Errorf("keyboard error: %w", err)
		}

		if err := v.handleKeyPress(char, key); err != nil {
			return err
		}
	}
	return nil
}

func (v *CLIViewer) displayCurrentPage() error {
	chapter, err := v.reader.GetChapter(v.currentPos.Chapter)
	if err != nil {
		return err
	}

	clearScreen()

	content := v.formatContent(chapter.Content)
	pages := v.paginateContent(content)

	if v.currentPos.Page >= len(pages) {
		v.currentPos.Page = len(pages) - 1
	}

	if v.currentPos.Page >= 0 && v.currentPos.Page < len(pages) {
		fmt.Println(pages[v.currentPos.Page])
	}

	footer := color.New(color.FgYellow)
	footer.Printf("\nPage %d of %d\n", v.currentPos.Page+1, len(pages))
	footer.Println("\nUse arrow keys to navigate (←/→ pages, ↑/↓ chapters)")
	footer.Println("Press 'h' for help, 'q' to quit")
	return nil
}

func (v *CLIViewer) nextPage() error {
	chapter, err := v.reader.GetChapter(v.currentPos.Chapter)
	if err != nil {
		return err
	}

	content := v.formatContent(chapter.Content)
	pages := v.paginateContent(content)

	if v.currentPos.Page < len(pages)-1 {
		v.currentPos.Page++
	} else {
		return v.nextChapter()
	}
	return nil
}

func (v *CLIViewer) previousPage() error {
	if v.currentPos.Page > 0 {
		v.currentPos.Page--
	} else {
		if err := v.previousChapter(); err != nil {
			return err
		}
		chapter, _ := v.reader.GetChapter(v.currentPos.Chapter)
		content := v.formatContent(chapter.Content)
		pages := v.paginateContent(content)
		v.currentPos.Page = len(pages) - 1
	}
	return nil
}

func (v *CLIViewer) nextChapter() error {
	if v.currentPos.Chapter < v.reader.GetTotalChapters()-1 {
		v.currentPos.Chapter++
		v.currentPos.Page = 0
	}
	return nil
}

func (v *CLIViewer) previousChapter() error {
	if v.currentPos.Chapter > 0 {
		v.currentPos.Chapter--
		v.currentPos.Page = 0
	}
	return nil
}

func (v *CLIViewer) goToChapter(num int) error {
	if num >= 0 && num < v.reader.GetTotalChapters() {
		v.currentPos.Chapter = num
		v.currentPos.Page = 0
	}
	return nil
}

func (v *CLIViewer) showTableOfContents() {
	clearScreen()
	toc := color.New(color.FgCyan)
	toc.Println("=== Table of Contents ===")
	fmt.Println()

	for i := 0; i < v.reader.GetTotalChapters(); i++ {
		chapter, err := v.reader.GetChapter(i)
		if err != nil {
			continue
		}
		fmt.Printf("%3d. %s\n", i+1, chapter.Title)
	}
	fmt.Println("\nPress any key to continue...")
	keyboard.GetKey()
}

// Update the CLIViewer struct

// Update the Start method
func (v *CLIViewer) Start(filePath string) error {
	if err := v.initializeInput(); err != nil {
		return err
	}

	defer v.cleanup()

	if err := keyboard.Open(); err != nil {
		return fmt.Errorf("failed to initialize keyboard: %w", err)
	}
	defer keyboard.Close()

	// Store the current file path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	v.currentFile = absPath

	if _, err := v.reader.Open(filePath); err != nil {
		return fmt.Errorf("failed to open book: %w", err)
	}
	defer v.reader.Close()

	// Load saved progress
	if err := v.loadProgress(); err != nil {
		// Log the error but continue with default position

		fmt.Fprintf(os.Stderr, "Warning: Failed to load progress: %v\n", err)
	}

	if err := v.showWelcomeScreen(); err != nil {
		return err
	}
	return v.eventLoop()
}

// Update the cleanup method
func (v *CLIViewer) cleanup() {
	// Save progress before closing
	if err := v.saveProgress(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save progress: %v\n", err)
	}

	if v.input != nil {
		v.input.Close()
	}
}

// Update handleKeyPress to save progress on exit
func (v *CLIViewer) handleKeyPress(char rune, key keyboard.Key) error {
	switch key {
	case keyboard.KeyArrowRight:
		return v.nextPage()
	case keyboard.KeyArrowLeft:
		return v.previousPage()
	case keyboard.KeyArrowUp:
		return v.previousChapter()
	case keyboard.KeyArrowDown:
		return v.nextChapter()
	case keyboard.KeyEsc:
		v.shouldExit = true
		// Save progress when exiting
		return v.saveProgress()
	}

	switch char {
	case 'q':
		v.shouldExit = true
		// Save progress when exiting
		return v.saveProgress()
	case 'h':
		v.showHelp()
	case 't':
		v.showTableOfContents()
	case 's':
		// Add manual save option
		if err := v.saveProgress(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save progress: %v\n", err)
		} else {
			fmt.Println("\nProgress saved successfully!")
			keyboard.GetKey() // Wait for key press
		}
	default:
		if num, err := strconv.Atoi(string(char)); err == nil {
			return v.goToChapter(num - 1)
		}
	}
	return nil
}

// Update showHelp to include save command
func (v *CLIViewer) showHelp() {
	clearScreen()
	help := color.New(color.FgCyan)
	help.Println("=== Help ===")
	fmt.Println()
	fmt.Println("Navigation:")
	fmt.Println("  → or Enter     - Next page")
	fmt.Println("  ←              - Previous page")
	fmt.Println("  ↓              - Next chapter")
	fmt.Println("  ↑              - Previous chapter")
	fmt.Println("  [number]       - Go to chapter number")
	fmt.Println()
	fmt.Println("Other commands:")
	fmt.Println("  h              - Show this help")
	fmt.Println("  t              - Show table of contents")
	fmt.Println("  s              - Save progress manually")
	fmt.Println("  q or ESC       - Exit the reader")
	fmt.Println("\nProgress is automatically saved when exiting")
	fmt.Println("\nPress any key to continue...")
	keyboard.GetKey()
}

// Update showWelcomeScreen to show loading progress
func (v *CLIViewer) showWelcomeScreen() error {

	metadata := v.reader.GetMetadata()
	clearScreen()

	title := color.New(color.FgCyan, color.Bold)
	title.Println("=== EPUB Reader ===")
	fmt.Println()
	fmt.Printf("Title: %s\n", metadata.Title)
	fmt.Printf("Author: %s\n", metadata.Author)
	fmt.Printf("Total Chapters: %d\n", v.reader.GetTotalChapters())

	if v.currentPos.Chapter > 0 || v.currentPos.Page > 0 {
		fmt.Printf("\nResuming from Chapter %d, Page %d\n",
			v.currentPos.Chapter+1, v.currentPos.Page+1)
	}

	fmt.Println("\nUse arrow keys to navigate")
	fmt.Println("Press 'h' for help, 'q' to quit")
	fmt.Println("Press any key to start reading...")

	_, _, err := keyboard.GetKey()
	return err
}
