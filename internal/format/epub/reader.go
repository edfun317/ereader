package epub

import (
	"archive/zip"
	"errors"

	"github.com/edfun317/ereader/internal/core"
)

type EPUBReader struct {
	file        *zip.ReadCloser
	book        *core.Book
	rootFile    string
	contentPath string
}

func NewEPUBReader() *EPUBReader {
	return &EPUBReader{}
}

func (r *EPUBReader) Open(path string) (*core.Book, error) {
	// Open the EPUB file (it's a ZIP file)
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	r.file = reader

	// Initialize book structure
	r.book = &core.Book{
		Chapters: make([]core.Chapter, 0),
	}

	// Read container.xml to find the root file
	if err := r.readContainer(); err != nil {
		return nil, err
	}

	// Read the package file (content.opf)
	if err := r.readPackage(); err != nil {
		return nil, err
	}

	return r.book, nil
}

func (r *EPUBReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func (r *EPUBReader) GetMetadata() core.BookMetadata {
	if r.book != nil {
		return r.book.Metadata
	}
	return core.BookMetadata{}
}

func (r *EPUBReader) GetChapter(index int) (*core.Chapter, error) {
	if r.book == nil {
		return nil, errors.New("book not opened")
	}
	if index < 0 || index >= len(r.book.Chapters) {
		return nil, errors.New("chapter index out of range")
	}
	return &r.book.Chapters[index], nil
}

func (r *EPUBReader) GetTotalChapters() int {
	if r.book == nil {
		return 0
	}
	return len(r.book.Chapters)
}
