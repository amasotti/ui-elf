package scanner

import (
	"regexp"
	"strings"

	"ui-elf/internal/types"
)

// VueParser parses Vue.js single-file components (.vue files)
// Extracts component usage from both template and script sections
type VueParser struct{}

// NewVueParser creates a new VueParser instance
func NewVueParser() *VueParser {
	return &VueParser{}
}

// SupportsFile checks if the file is a .vue file
func (p *VueParser) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".vue")
}

// Parse extracts component matches from Vue file content
// Handles both template syntax and JSX in script sections
func (p *VueParser) Parse(fileContent string, filePath string) ([]types.ComponentMatch, error) {
	var matches []types.ComponentMatch

	// Extract template section
	templateContent, templateStartLine := extractTemplateSection(fileContent)
	if templateContent != "" {
		templateMatches := parseTemplateComponents(templateContent, filePath, templateStartLine)
		matches = append(matches, templateMatches...)
	}

	// Extract script section and look for JSX
	scriptContent, scriptStartLine := extractScriptSection(fileContent)
	if scriptContent != "" {
		jsxMatches := parseJSXComponents(scriptContent, filePath, scriptStartLine)
		matches = append(matches, jsxMatches...)
	}

	return matches, nil
}

// extractTemplateSection extracts the content within <template> tags
// Returns the template content and the line number where the template starts
func extractTemplateSection(content string) (string, int) {
	// Match <template> or <template lang="..."> tags
	templateRegex := regexp.MustCompile(`(?s)<template[^>]*>(.*?)</template>`)
	match := templateRegex.FindStringSubmatchIndex(content)

	if len(match) < 4 {
		return "", 0
	}

	// Extract the template content (first capture group)
	templateContent := content[match[2]:match[3]]

	// Calculate the starting line number
	startLine := strings.Count(content[:match[2]], "\n") + 1

	return templateContent, startLine
}

// extractScriptSection extracts the content within <script> tags
// Returns the script content and the line number where the script starts
func extractScriptSection(content string) (string, int) {
	// Match <script> or <script lang="..."> or <script setup> tags
	scriptRegex := regexp.MustCompile(`(?s)<script[^>]*>(.*?)</script>`)
	match := scriptRegex.FindStringSubmatchIndex(content)

	if len(match) < 4 {
		return "", 0
	}

	// Extract the script content (first capture group)
	scriptContent := content[match[2]:match[3]]

	// Calculate the starting line number
	startLine := strings.Count(content[:match[2]], "\n") + 1

	return scriptContent, startLine
}

// parseTemplateComponents extracts component usage from template content
// Matches both self-closing and paired tags: <ComponentName /> and <ComponentName>
func parseTemplateComponents(templateContent string, filePath string, baseLineNumber int) []types.ComponentMatch {
	var matches []types.ComponentMatch

	// Regex to match opening tags - <tagname followed by whitespace, >, /, or end of line
	// This handles multi-line tags where attributes span multiple lines
	componentRegex := regexp.MustCompile(`<([A-Za-z][A-Za-z0-9-]*)(?:[\s>/]|$)`)

	lines := strings.Split(templateContent, "\n")
	seenComponents := make(map[string]map[int]bool) // Track component:line to avoid duplicates

	for lineIdx, line := range lines {
		componentMatches := componentRegex.FindAllStringSubmatch(line, -1)

		for _, match := range componentMatches {
			if len(match) >= 2 {
				componentName := match[1]

				// Skip HTML tags (lowercase only, no hyphens or uppercase)
				if isHTMLTag(componentName) {
					continue
				}

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

// parseJSXComponents extracts component usage from JSX syntax in script sections
// Handles JSX elements like <Component /> or <Component>
func parseJSXComponents(scriptContent string, filePath string, baseLineNumber int) []types.ComponentMatch {
	var matches []types.ComponentMatch

	// Regex to match JSX component tags
	// JSX components must start with uppercase letter
	componentRegex := regexp.MustCompile(`<([A-Z][A-Za-z0-9]*)(?:[\s>/]|$)`)

	lines := strings.Split(scriptContent, "\n")
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

// isHTMLTag checks if a tag name is a standard HTML element
// Returns true for common HTML tags that should be ignored
func isHTMLTag(tagName string) bool {
	// Common HTML tags (lowercase only)
	htmlTags := map[string]bool{
		"div": true, "span": true, "p": true, "a": true, "img": true,
		"ul": true, "ol": true, "li": true, "table": true, "tr": true,
		"td": true, "th": true, "thead": true, "tbody": true, "tfoot": true,
		"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
		"header": true, "footer": true, "nav": true, "section": true, "article": true,
		"aside": true, "main": true, "input": true, "textarea": true, "select": true,
		"option": true, "label": true, "fieldset": true, "legend": true,
		"strong": true, "em": true, "b": true, "i": true, "u": true,
		"br": true, "hr": true, "pre": true, "code": true, "blockquote": true,
		"iframe": true, "video": true, "audio": true, "canvas": true, "svg": true,
		"path": true, "circle": true, "rect": true, "line": true, "polygon": true,
		"template": true, "slot": true, "script": true, "style": true, "link": true,
		"meta": true, "title": true, "head": true, "body": true, "html": true,
		"button": true, "form": true, "dialog": true,
	}

	// Check if it's a standard HTML tag (must be all lowercase and in the map)
	lowerTag := strings.ToLower(tagName)
	return lowerTag == tagName && htmlTags[lowerTag]
}
