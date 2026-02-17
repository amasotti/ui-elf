package types

// ComponentMatch represents a single component found in the codebase
type ComponentMatch struct {
	FilePath      string `json:"filePath"`      // Relative path to the file
	Line          int    `json:"line"`          // Line number where component appears
	ComponentName string `json:"componentName"` // Actual component name (e.g., "q-form")
	ComponentType string `json:"componentType"` // Normalized type (e.g., "form")
}

// ScanResult contains aggregated results from scanning the codebase
type ScanResult struct {
	Matches       []ComponentMatch `json:"matches"`
	TotalCount    int              `json:"totalCount"`
	ScanTimeMs    int64            `json:"scanTimeMs"`
	ComponentType string           `json:"componentType"`
	ScannedFiles  int              `json:"scannedFiles"`
}

// CLIOptions holds parsed command-line arguments
type CLIOptions struct {
	ComponentType string
	Directory     string
	Filter        []string
	OutputFormat  string // "terminal", "json", or "both"
}

// FileFilter defines criteria for filtering files during discovery
type FileFilter struct {
	ExcludePatterns   []string
	IncludeDirectories []string
	FileExtensions    []string
}
