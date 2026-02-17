// Package output handles the formatting and display of scan results.
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"ui-elf/internal/types"
)

// OutputFormatter handles formatting and displaying scan results
type OutputFormatter struct{}

// NewOutputFormatter creates a new output formatter
func NewOutputFormatter() *OutputFormatter {
	return &OutputFormatter{}
}

// FormatTerminal formats the scan result for terminal display
// Shows file paths, counts, and scan time
func (f *OutputFormatter) FormatTerminal(result *types.ScanResult) string {
	var sb strings.Builder

	// Header
	fmt.Fprintf(&sb, "\nComponent Finder Results - %s\n", result.ComponentType)
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	// File paths
	if len(result.Matches) == 0 {
		sb.WriteString("No components found.\n")
	} else {
		sb.WriteString("Found components in:\n\n")
		for _, match := range result.Matches {
			fmt.Fprintf(&sb, "  %s (line %d): %s\n",
				match.FilePath, match.Line, match.ComponentName)
		}
	}

	// Summary
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("-", 50))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "Total components found: %d\n", result.TotalCount)
	fmt.Fprintf(&sb, "Files scanned: %d\n", result.ScannedFiles)
	fmt.Fprintf(&sb, "Scan time: %dms\n", result.ScanTimeMs)

	return sb.String()
}

// FormatJSON formats the scan result as JSON
// Returns a JSON string with all result data
func (f *OutputFormatter) FormatJSON(result *types.ScanResult) (string, error) {
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// Write outputs the scan result according to the specified options
// Supports terminal, JSON file output, or both
func (f *OutputFormatter) Write(result *types.ScanResult, format string, outputPath string) error {
	switch format {
	case "terminal":
		fmt.Print(f.FormatTerminal(result))

	case "json":
		jsonStr, err := f.FormatJSON(result)
		if err != nil {
			return err
		}

		if outputPath == "" {
			outputPath = "ui-elf-results.json"
		}

		if err := os.WriteFile(outputPath, []byte(jsonStr), 0644); err != nil {
			return fmt.Errorf("failed to write JSON file: %w", err)
		}

		fmt.Printf("Results written to %s\n", outputPath)

	case "both":
		// Display terminal output
		fmt.Print(f.FormatTerminal(result))

		// Write JSON file
		jsonStr, err := f.FormatJSON(result)
		if err != nil {
			return err
		}

		if outputPath == "" {
			outputPath = "ui-elf-results.json"
		}

		if err := os.WriteFile(outputPath, []byte(jsonStr), 0644); err != nil {
			return fmt.Errorf("failed to write JSON file: %w", err)
		}

		fmt.Printf("\nResults also written to %s\n", outputPath)

	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	return nil
}
