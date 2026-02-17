// Package cli implements the command-line interface controller and logic.
package cli

import (
	"fmt"
	"os"

	"ui-elf/internal/discovery"
	"ui-elf/internal/output"
	"ui-elf/internal/registry"
	"ui-elf/internal/scanner"
	"ui-elf/internal/types"

	"github.com/spf13/cobra"
)

// Controller orchestrates the CLI operations
type Controller struct {
	rootCmd *cobra.Command
}

// NewController creates a new CLI controller with cobra configuration
func NewController() *Controller {
	c := &Controller{}
	c.setupRootCommand()
	return c
}

// setupRootCommand configures the root cobra command with flags and help text
func (c *Controller) setupRootCommand() {
	c.rootCmd = &cobra.Command{
		Use:   "ui-elf [flags]",
		Short: "Scan Vue.js and React codebases for specific component types",
		Long: `UI Elf scans your codebase to locate specific component types
(forms, buttons, dialogs, and custom components) in Vue.js and React projects.

The tool helps development teams audit their frontend applications by identifying
where components are used and providing usage statistics.`,
		Example: `  # Scan for forms in current directory
  ui-elf --component-type form --directory .

  # Scan for buttons in src directory with JSON output
  ui-elf --component-type button --directory ./src --output json

  # Scan for custom component with directory filter
  ui-elf --component-type custom --directory . --filter src/components,src/views

  # Scan for dialogs with both terminal and JSON output
  ui-elf --component-type dialog --directory . --output both`,
		RunE: c.run,
	}

	// Define flags
	c.rootCmd.Flags().StringP("component-type", "t", "", "Component type to search for (form, button, dialog, custom) [required]")
	c.rootCmd.Flags().StringP("directory", "d", ".", "Directory to scan (default: current directory)")
	c.rootCmd.Flags().StringSliceP("filter", "f", []string{}, "Comma-separated list of directories to include (e.g., src/components,src/views)")
	c.rootCmd.Flags().StringP("output", "o", "terminal", "Output format: terminal, json, or both (default: terminal)")

	// Mark required flags
	if err := c.rootCmd.MarkFlagRequired("component-type"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking flag required: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main CLI logic
func (c *Controller) run(cmd *cobra.Command, args []string) error {
	// Parse flags into CLIOptions
	options, err := c.parseFlags(cmd)
	if err != nil {
		return err
	}

	// Validate options
	if err := c.validateOptions(options); err != nil {
		return err
	}

	// Execute the scan
	result, err := c.executeScan(options)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Format and display output
	if err := c.displayOutput(result, options); err != nil {
		return fmt.Errorf("failed to display output: %w", err)
	}

	return nil
}

// parseFlags extracts flag values into CLIOptions struct
func (c *Controller) parseFlags(cmd *cobra.Command) (*types.CLIOptions, error) {
	componentType, err := cmd.Flags().GetString("component-type")
	if err != nil {
		return nil, fmt.Errorf("failed to parse component-type flag: %w", err)
	}

	directory, err := cmd.Flags().GetString("directory")
	if err != nil {
		return nil, fmt.Errorf("failed to parse directory flag: %w", err)
	}

	filter, err := cmd.Flags().GetStringSlice("filter")
	if err != nil {
		return nil, fmt.Errorf("failed to parse filter flag: %w", err)
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, fmt.Errorf("failed to parse output flag: %w", err)
	}

	return &types.CLIOptions{
		ComponentType: componentType,
		Directory:     directory,
		Filter:        filter,
		OutputFormat:  output,
	}, nil
}

// validateOptions validates the parsed CLI options
func (c *Controller) validateOptions(options *types.CLIOptions) error {
	// Validate component type
	validTypes := map[string]bool{
		"form":   true,
		"button": true,
		"dialog": true,
		"custom": true,
	}
	if !validTypes[options.ComponentType] {
		return fmt.Errorf("invalid component type '%s': must be one of: form, button, dialog, custom", options.ComponentType)
	}

	// Validate output format
	validOutputs := map[string]bool{
		"terminal": true,
		"json":     true,
		"both":     true,
	}
	if !validOutputs[options.OutputFormat] {
		return fmt.Errorf("invalid output format '%s': must be one of: terminal, json, both", options.OutputFormat)
	}

	// Validate directory exists
	if _, err := os.Stat(options.Directory); os.IsNotExist(err) {
		return fmt.Errorf("directory not found: %s", options.Directory)
	}

	return nil
}

// Execute runs the CLI controller
func (c *Controller) Execute() error {
	return c.rootCmd.Execute()
}

// executeScan performs the component scanning process
func (c *Controller) executeScan(options *types.CLIOptions) (*types.ScanResult, error) {
	// Import required packages at the top of the file
	// Create file discovery service
	discoveryService := discovery.NewFileDiscoveryService()

	// Build file filter
	filter := types.FileFilter{
		ExcludePatterns:    []string{"node_modules", "test", "tests", "__tests__", ".test.", ".spec."},
		IncludeDirectories: options.Filter,
		FileExtensions:     []string{".vue", ".jsx", ".tsx"},
	}

	// Discover files
	files, err := discoveryService.DiscoverFiles(options.Directory, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to discover files: %w", err)
	}

	// Check if any files were found
	if len(files) == 0 {
		return &types.ScanResult{
			Matches:       []types.ComponentMatch{},
			TotalCount:    0,
			ScanTimeMs:    0,
			ComponentType: options.ComponentType,
			ScannedFiles:  0,
		}, nil
	}

	// Create component registry
	registry := registry.NewComponentMappingRegistry()

	// Create parsers
	parsers := []scanner.ComponentParser{
		scanner.NewVueParser(),
		scanner.NewReactParser(),
	}

	// Create scanner
	componentScanner := scanner.NewComponentScanner(parsers, registry)

	// Execute scan
	result, err := componentScanner.Scan(files, options.ComponentType)
	if err != nil {
		return nil, fmt.Errorf("scan execution failed: %w", err)
	}

	return result, nil
}

// displayOutput formats and displays the scan results
func (c *Controller) displayOutput(result *types.ScanResult, options *types.CLIOptions) error {
	formatter := output.NewOutputFormatter()

	// Determine output path for JSON (empty string will use default)
	outputPath := ""

	// Write output according to format
	if err := formatter.Write(result, options.OutputFormat, outputPath); err != nil {
		return err
	}

	return nil
}
