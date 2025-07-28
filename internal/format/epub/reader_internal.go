package epub

import (
	"encoding/xml"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/edfun317/ereader/internal/core"
)

// Internal helper methods
func (r *EPUBReader) readContainer() error {
	containerFile, err := r.findFile("META-INF/container.xml")
	if err != nil {
		return err
	}
	defer containerFile.Close()

	var container container
	if err := xml.NewDecoder(containerFile).Decode(&container); err != nil {
		return err
	}

	if len(container.RootFiles) == 0 {
		return errors.New("no rootfile found in container.xml")
	}

	r.rootFile = container.RootFiles[0].FullPath
	r.contentPath = filepath.Dir(r.rootFile)
	return nil
}

func (r *EPUBReader) readPackage() error {
	packageFile, err := r.findFile(r.rootFile)
	if err != nil {
		return err
	}
	defer packageFile.Close()

	var pkg Package
	if err := xml.NewDecoder(packageFile).Decode(&pkg); err != nil {
		return err
	}

	// Set metadata
	r.book.Metadata = core.BookMetadata{
		Title:     pkg.Metadata.Title,
		Author:    pkg.Metadata.Creator,
		Publisher: pkg.Metadata.Publisher,
		Language:  pkg.Metadata.Language,
	}

	// Build a map of manifest items for quick lookup
	manifestItems := make(map[string]Item)
	for _, item := range pkg.Manifest.Items {
		manifestItems[item.ID] = item
	}

	// Process chapters in spine order
	for i, itemRef := range pkg.Spine.ItemRefs {
		item, ok := manifestItems[itemRef.IDRef]
		if !ok {
			continue
		}

		if item.MediaType == "application/xhtml+xml" ||
			item.MediaType == "application/x-dtbook+xml" {
			chapter, err := r.readChapter(item, i)
			if err != nil {
				continue // Skip chapters with errors
			}
			r.book.Chapters = append(r.book.Chapters, *chapter)
		}
	}

	return nil
}

func (r *EPUBReader) readChapter(item Item, index int) (*core.Chapter, error) {
	chapterPath := filepath.Join(r.contentPath, item.Href)
	chapterFile, err := r.findFile(chapterPath)
	if err != nil {
		return nil, err
	}
	defer chapterFile.Close()

	content, err := io.ReadAll(chapterFile)
	if err != nil {
		return nil, err
	}

	return &core.Chapter{
		Index:   index,
		Title:   extractTitle(content),
		Content: string(content),
	}, nil
}

func (r *EPUBReader) findFile(name string) (io.ReadCloser, error) {
	for _, f := range r.file.File {
		if strings.EqualFold(f.Name, name) {
			return f.Open()
		}
	}
	return nil, errors.New("file not found in EPUB: " + name)
}

// Helper function to extract title from chapter content
func extractTitle(content []byte) string {
	titleStart := strings.Index(string(content), "<title>")
	if titleStart == -1 {
		return "Untitled"
	}
	titleStart += 7 // len("<title>")

	titleEnd := strings.Index(string(content), "</title>")
	if titleEnd == -1 {
		return "Untitled"
	}
	return strings.TrimSpace(string(content[titleStart:titleEnd]))
}
