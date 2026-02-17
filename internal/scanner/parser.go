package scanner

import "ui-elf/internal/types"

// ComponentParser defines the interface for parsing component files
// Implementations should handle specific file types (Vue, React, etc.)
type ComponentParser interface {
	// Parse extracts component matches from the given file content
	// Returns a slice of ComponentMatch instances found in the file
	// Requirements: 2.1 (Vue parsing), 2.2 (React parsing)
	Parse(fileContent string, filePath string) ([]types.ComponentMatch, error)

	// SupportsFile determines if this parser can handle the given file
	// Returns true if the parser supports the file extension/type
	// Requirements: 2.1 (Vue files), 2.2 (React files)
	SupportsFile(filePath string) bool
}
