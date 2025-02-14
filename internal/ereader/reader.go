// Package ebookreader provides functionality for reading and parsing EPUB files
package ebookreader

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// XML structure definitions for EPUB parsing

// Container represents the structure of container.xml
type Container struct {
	XMLName   xml.Name   `xml:"container"`
	Rootfiles []Rootfile `xml:"rootfiles>rootfile"`
}

// Rootfile represents a rootfile entry in container.xml
type Rootfile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

// OPF represents the structure of the OPF (Open Packaging Format) file
type OPF struct {
	XMLName  xml.Name `xml:"package"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

// Manifest contains the list of resources in the EPUB
type Manifest struct {
	Items []Item `xml:"item"`
}

// Item represents a resource item in the manifest
type Item struct {
	ID   string `xml:"id,attr"`
	Href string `xml:"href,attr"`
}

// Spine defines the reading order of content documents
type Spine struct {
	ItemRefs []ItemRef `xml:"itemref"`
}

// ItemRef represents a reference to a manifest item in the spine
type ItemRef struct {
	IDRef string `xml:"idref,attr"`
}

// EPUBReader implements the Reader interface for EPUB format
type EPUBReader struct {
	file     *zip.ReadCloser
	spine    []string          // Ordered list of content files
	manifest map[string]string // Resource mapping (ID to file path)
	current  int               // Current reading position in spine
}

// Open initializes the EPUB reader with the specified file
// It parses the EPUB structure and prepares for content extraction
func (r *EPUBReader) Open(filePath string) error {
	// Open EPUB file (ZIP format)
	file, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	r.file = file

	// Find and parse the OPF file path from container.xml
	opfPath, err := findOPFFile(file)
	if err != nil {
		return err
	}

	// Open and parse the OPF file
	opfFile, err := file.Open(opfPath)
	if err != nil {
		return err
	}
	defer opfFile.Close()

	// Decode OPF file content
	var opf OPF
	decoder := xml.NewDecoder(opfFile)
	if err := decoder.Decode(&opf); err != nil {
		return err
	}

	// Initialize spine and manifest
	r.spine = make([]string, 0)
	r.manifest = make(map[string]string)

	// Build manifest mapping
	for _, item := range opf.Manifest.Items {
		r.manifest[item.ID] = item.Href
	}

	// Build spine order
	for _, itemRef := range opf.Spine.ItemRefs {
		r.spine = append(r.spine, itemRef.IDRef)
	}
	r.current = 0
	return nil
}

// Next retrieves the next content section from the EPUB
// Returns the text content and any error encountered
func (r *EPUBReader) Next() (string, error) {
	if r.current >= len(r.spine) {
		return "", io.EOF // All content has been read
	}

	itemID := r.spine[r.current]
	itemPath, ok := r.manifest[itemID]
	if !ok {
		return "", errors.New("invalid spine item")
	}

	// Open the corresponding HTML file
	itemFile, err := r.file.Open(filepath.Join("OEBPS", itemPath))
	if err != nil {
		return "", err
	}
	defer itemFile.Close()

	// Read HTML content
	content, err := io.ReadAll(itemFile)
	if err != nil {
		return "", err
	}

	// Extract plain text from HTML
	textContent, err := extractTextFromHTML(content)
	if err != nil {
		return "", err
	}

	r.current++
	return textContent, nil
}

// Close releases resources associated with the EPUB reader
func (r *EPUBReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// SupportsFormat checks if the specified format is supported
func (r *EPUBReader) SupportsFormat(format string) bool {
	return strings.ToLower(format) == "epub"
}

// findOPFFile locates the OPF(Open Packaging Format) file path within the EPUB archive
func findOPFFile(file *zip.ReadCloser) (string, error) {

	// Open container.xml from META-INF directory
	containerFile, err := file.Open("META-INF/container.xml")
	if err != nil {
		return "", err
	}
	defer containerFile.Close()

	// Parse container.xml content
	var container Container
	decoder := xml.NewDecoder(containerFile)
	if err := decoder.Decode(&container); err != nil {
		return "", err
	}

	// Find and return the OPF file path
	for _, rootfile := range container.Rootfiles {
		if rootfile.MediaType == "application/oebps-package+xml" {
			return rootfile.FullPath, nil
		}
	}
	return "", errors.New("OPF file not found")
}

// extractTextFromHTML extracts plain text content from HTML
func extractTextFromHTML(htmlContent []byte) (string, error) {

	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var textContent strings.Builder
	var traverse func(*html.Node)

	// Define recursive traversal function
	traverse = func(n *html.Node) {
		if n.Type == html.TextNode {
			textContent.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return textContent.String(), nil
}
