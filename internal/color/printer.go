package color

import (
	"strings"

	"github.com/fatih/color"
)

// Printer handles colored text output
type Printer struct {
	scheme    ColorScheme
	textColor color.Attribute
	bgColor   color.Attribute
}

// NewPrinter creates a new Printer instance
func NewPrinter(schemeName string) *Printer {
	scheme, ok := PredefinedSchemes[schemeName]
	if !ok {
		scheme = PredefinedSchemes["default"]
	}

	return &Printer{
		scheme:    scheme,
		textColor: ColorMap[strings.ToLower(scheme.TextColor)],
		bgColor:   BgColorMap[strings.ToLower(scheme.BgColor)],
	}
}

// Print prints colored text
func (p *Printer) Print(text string) {
	c := color.New(p.textColor, p.bgColor)
	c.Print(text)
}

// Println prints colored text with newline
func (p *Printer) Println(text string) {
	c := color.New(p.textColor, p.bgColor)
	c.Println(text)

}

func (p *Printer) Printf(format string, a ...interface{}) {

	c := color.New(p.textColor, p.bgColor)
	c.Printf(format, a...)
}

// GetCurrentScheme returns the current color scheme
func (p *Printer) GetCurrentScheme() ColorScheme {
	return p.scheme
}
