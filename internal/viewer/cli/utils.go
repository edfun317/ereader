package cli

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/net/html"
)

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// formatContent
func (v *CLIViewer) formatContent(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return content
	}

	var buf bytes.Buffer
	var extract func(*html.Node)

	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
			buf.WriteString(" ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
		if n.Type == html.ElementNode && (n.Data == "p" || n.Data == "div") {
			buf.WriteString("\n\n")
		}
	}

	extract(doc)
	return strings.TrimSpace(buf.String())
}

func (v *CLIViewer) paginateContent(content string) []string {
	const maxWidth = 80
	lines := make([]string, 0)
	paragraphs := strings.Split(content, "\n\n")

	for _, paragraph := range paragraphs {
		if strings.TrimSpace(paragraph) == "" {
			continue
		}

		words := strings.Fields(paragraph)
		if len(words) == 0 {
			continue
		}

		var currentLine strings.Builder
		currentLine.WriteString(words[0])
		lineLength := len(words[0])

		for i := 1; i < len(words); i++ {
			word := words[i]
			wordLen := len(word)

			if lineLength+wordLen+1 > maxWidth {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
				lineLength = wordLen
			} else {
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
				lineLength += wordLen + 1
			}
		}

		if currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		}
		lines = append(lines, "") // Add blank line between paragraphs
	}

	// Group lines into pages
	var pages []string
	var currentPage strings.Builder
	lineCount := 0

	for _, line := range lines {
		if lineCount >= v.pageSize {
			pages = append(pages, strings.TrimSpace(currentPage.String()))
			currentPage.Reset()
			lineCount = 0
		}

		currentPage.WriteString(line)
		currentPage.WriteString("\n")
		lineCount++
	}

	if currentPage.Len() > 0 {
		pages = append(pages, strings.TrimSpace(currentPage.String()))
	}

	return pages
}
