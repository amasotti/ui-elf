package scanner

import (
	"regexp"
	"strings"

	"ui-elf/internal/types"
)

// ReactParser parses React component files (.jsx and .tsx files)
// Extracts component usage from JSX elements
type ReactParser struct{}

// NewReactParser creates a new ReactParser instance
func NewReactParser() *ReactParser {
	return &ReactParser{}
}

// SupportsFile checks if the file is a .jsx or .tsx file
func (p *ReactParser) SupportsFile(filePath string) bool {
	lowerPath := strings.ToLower(filePath)
	return strings.HasSuffix(lowerPath, ".jsx") || strings.HasSuffix(lowerPath, ".tsx")
}

// Parse extracts component matches from React file content
// Handles JSX syntax in both .jsx and .tsx files
func (p *ReactParser) Parse(fileContent string, filePath string) ([]types.ComponentMatch, error) {
	return parseReactJSXComponents(fileContent, filePath, 1), nil
}

// parseReactJSXComponents extracts component usage from JSX syntax
// Handles JSX elements like <Component /> or <Component>
func parseReactJSXComponents(content string, filePath string, baseLineNumber int) []types.ComponentMatch {
	var matches []types.ComponentMatch
	
	// Regex to match JSX component tags
	// JSX components must start with uppercase letter
	// Matches: <ComponentName followed by whitespace, >, /, or end of line
	componentRegex := regexp.MustCompile(`<([A-Z][A-Za-z0-9]*)(?:[\s>/]|$)`)
	
	lines := strings.Split(content, "\n")
	seenComponents := make(map[string]map[int]bool) // Track component:line to avoid duplicates
	
	for lineIdx, line := range lines {
		componentMatches := componentRegex.FindAllStringSubmatch(line, -1)
		
		for _, match := range componentMatches {
			if len(match) >= 2 {
				componentName := match[1]
				
				// Skip if we've already seen this component on this line
				if seenComponents[componentName] == nil {
					seenComponents[componentName] = make(map[int]bool)
				}
				if seenComponents[componentName][lineIdx] {
					continue
				}
				seenComponents[componentName][lineIdx] = true
				
				matches = append(matches, types.ComponentMatch{
					FilePath:      filePath,
					Line:          baseLineNumber + lineIdx,
					ComponentName: componentName,
					ComponentType: "", // Will be set by scanner based on registry
				})
			}
		}
	}
	
	return matches
}
