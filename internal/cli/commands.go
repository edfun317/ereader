package cli

import (
	"fmt"
	"io"

	"github.com/edfun317/ereader/internal/color"
	"github.com/edfun317/ereader/internal/reader"
	"github.com/spf13/cobra"
)

var (
	schemeName string
	rootCmd    = &cobra.Command{
		Use:   "ereader",
		Short: "A text reader with color support",
		Long:  `A command line text reader that supports colored output and various text formats.`,
	}
)

// InitCommands initializes all CLI commands
func InitCommands() *cobra.Command {

	// Add persistent flags
	rootCmd.PersistentFlags().StringVarP(&schemeName, "scheme", "s", "default", "Use predefined color scheme")

	// Add commands
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(schemesCmd)

	return rootCmd
}

var readCmd = &cobra.Command{
	Use:   "read [filepath]",
	Short: "Read a text file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := args[0]

		// Create file reader
		fr, err := reader.NewFileReader(filepath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer fr.Close()

		// Create color printer
		printer := color.NewPrinter(schemeName)

		// Read and print file content
		for {
			line, err := fr.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			printer.Println(line)
		}

		return nil
	},
}

var schemesCmd = &cobra.Command{
	Use:   "schemes",
	Short: "List available color schemes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available color schemes:")
		for _, scheme := range color.PredefinedSchemes {
			printer := color.NewPrinter(scheme.Name)
			printer.Printf("%-12s: %s\n", scheme.Name, scheme.Description)
		}
	},
}
