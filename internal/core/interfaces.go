// internal/core/interfaces.go
package core

// BookReader defines the interface for reading different ebook formats
type BookReader interface {
	Open(path string) (*Book, error)
	Close() error
	GetMetadata() BookMetadata
	GetChapter(index int) (*Chapter, error)
	GetTotalChapters() int
}

// Viewer defines the interface for different viewing methods
type Viewer interface {
	Initialize() error
	Display(content *Chapter) error
	ShowMetadata(metadata BookMetadata) error
	HandleUserInput() error
	Cleanup() error
}
