package reader

import (
	"bufio"
	"io"
	"os"
)

// FileReader handles file reading operations
type FileReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

// NewFileReader creates a new FileReader instance
func NewFileReader(filepath string) (*FileReader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

// ReadLine reads next line from file
func (r *FileReader) ReadLine() (string, error) {
	if r.scanner.Scan() {
		return r.scanner.Text(), nil
	}

	if err := r.scanner.Err(); err != nil {
		return "", err
	}

	return "", io.EOF
}

// Close closes the file
func (r *FileReader) Close() error {
	return r.file.Close()
}
