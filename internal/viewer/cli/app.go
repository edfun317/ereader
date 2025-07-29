// internal/viewer/cli/app.go
package cli

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/edfun317/ereader/internal/core"
)

type (
	CLIViewer struct {
		reader      core.BookReader
		scanner     *bufio.Scanner
		currentPos  CurrentPos
		pageSize    int
		shouldExit  bool
		input       *os.File
		currentFile string // Add this field to store current file path
	}
	CurrentPos struct {
		Chapter int `json:"chapter"`
		Page    int `json:"page"`
	}
)

func NewCLIViewer(reader core.BookReader) *CLIViewer {

	return &CLIViewer{
		reader:   reader,
		pageSize: 20, // Default lines per page
	}
}

// Initialize input handling
func (v *CLIViewer) initializeInput() error {
	// Open /dev/tty for Unix systems or CONIN$ for Windows
	var err error
	if runtime.GOOS == "windows" {
		v.input, err = os.Open("CONIN$")
	} else {
		v.input, err = os.Open("/dev/tty")
	}
	if err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}
	v.scanner = bufio.NewScanner(v.input)
	return nil
}
