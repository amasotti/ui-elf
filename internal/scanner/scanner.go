package scanner

import (
	"os"
	"sync"
	"time"

	"ui-elf/internal/registry"
	"ui-elf/internal/types"
)

// ComponentScanner coordinates the scanning process across multiple files
// Uses concurrent processing with goroutines for performance
type ComponentScanner struct {
	parsers  []ComponentParser
	registry *registry.ComponentMappingRegistry
}

// NewComponentScanner creates a new scanner with the given parsers
func NewComponentScanner(parsers []ComponentParser, reg *registry.ComponentMappingRegistry) *ComponentScanner {
	return &ComponentScanner{
		parsers:  parsers,
		registry: reg,
	}
}

// Scan processes all files concurrently and returns aggregated results
// Filters matches by component type using the registry
func (s *ComponentScanner) Scan(files []string, componentType string) (*types.ScanResult, error) {
	startTime := time.Now()
	
	// Channel to collect matches from all goroutines
	matchChan := make(chan []types.ComponentMatch, len(files))
	
	// WaitGroup to track completion of all goroutines
	var wg sync.WaitGroup
	
	// Process files concurrently
	for _, filePath := range files {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			
			// Find appropriate parser for this file
			var parser ComponentParser
			for _, p := range s.parsers {
				if p.SupportsFile(path) {
					parser = p
					break
				}
			}
			
			if parser == nil {
				// No parser supports this file, skip it
				matchChan <- nil
				return
			}
			
			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				// Log error but continue with other files
				// In production, we'd use a proper logger
				matchChan <- nil
				return
			}
			
			// Parse the file
			matches, err := parser.Parse(string(content), path)
			if err != nil {
				// Log error but continue with other files
				matchChan <- nil
				return
			}
			
			// Filter matches by component type
			filteredMatches := s.filterByComponentType(matches, componentType)
			matchChan <- filteredMatches
		}(filePath)
	}
	
	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(matchChan)
	}()
	
	// Collect all matches
	var allMatches []types.ComponentMatch
	for matches := range matchChan {
		if matches != nil {
			allMatches = append(allMatches, matches...)
		}
	}
	
	// Calculate scan time
	scanTime := time.Since(startTime)
	
	// Build result
	result := &types.ScanResult{
		Matches:       allMatches,
		TotalCount:    len(allMatches),
		ScanTimeMs:    scanTime.Milliseconds(),
		ComponentType: componentType,
		ScannedFiles:  len(files),
	}
	
	return result, nil
}

// filterByComponentType filters matches to only include those matching the component type
// Sets the ComponentType field on matching components
func (s *ComponentScanner) filterByComponentType(matches []types.ComponentMatch, componentType string) []types.ComponentMatch {
	var filtered []types.ComponentMatch
	
	for _, match := range matches {
		if s.registry.MatchesComponentType(match.ComponentName, componentType) {
			// Set the component type on the match
			match.ComponentType = componentType
			filtered = append(filtered, match)
		}
	}
	
	return filtered
}
