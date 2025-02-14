// internal/core/models.go
package core

type Book struct {
	Metadata BookMetadata
	Chapters []Chapter
}

type BookMetadata struct {
	Title     string
	Author    string
	Publisher string
	Language  string
}

type Chapter struct {
	Index   int
	Title   string
	Content string
}
