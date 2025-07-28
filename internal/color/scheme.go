package color

import "github.com/fatih/color"

// ColorScheme represents a predefined color combination
type ColorScheme struct {
	Name        string
	TextColor   string
	BgColor     string
	Description string
}

// Available colors mapping
var (
	ColorMap = map[string]color.Attribute{
		"red":     color.FgRed,
		"green":   color.FgGreen,
		"blue":    color.FgBlue,
		"yellow":  color.FgYellow,
		"magenta": color.FgMagenta,
		"cyan":    color.FgCyan,
		"white":   color.FgWhite,
		"black":   color.FgBlack,
	}

	BgColorMap = map[string]color.Attribute{
		"red":     color.BgRed,
		"green":   color.BgGreen,
		"blue":    color.BgBlue,
		"yellow":  color.BgYellow,
		"magenta": color.BgMagenta,
		"cyan":    color.BgCyan,
		"white":   color.BgWhite,
		"black":   color.BgBlack,
	}

	// PredefinedSchemes contains all available color schemes
	PredefinedSchemes = map[string]ColorScheme{
		"default": {
			Name:        "default",
			TextColor:   "white",
			BgColor:     "black",
			Description: "Classic terminal style, suitable for long-term use",
		},
		"paper": {
			Name:        "paper",
			TextColor:   "black",
			BgColor:     "white",
			Description: "Paper-like effect, suitable for reading long texts",
		},
		"night": {
			Name:        "night",
			TextColor:   "cyan",
			BgColor:     "black",
			Description: "Night mode, reduces eye strain",
		},
	}
)
