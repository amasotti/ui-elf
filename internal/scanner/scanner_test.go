package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"component-finder-cli/internal/registry"
	"component-finder-cli/internal/types"
)

func TestComponentScanner_Scan(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	
	// Create a Vue file with a form component
	vueFile := filepath.Join(tempDir, "test.vue")
	vueContent := `<template>
  <div>
    <q-form>
      <input type="text" />
    </q-form>
  </div>
</template>`
	err := os.WriteFile(vueFile, []byte(vueContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test Vue file: %v", err)
	}
	
	// Create a React file with a button component
	reactFile := filepath.Join(tempDir, "test.jsx")
	reactContent := `import React from 'react';

function MyComponent() {
  return (
    <div>
      <Button onClick={handleClick}>Click me</Button>
    </div>
  );
}`
	err = os.WriteFile(reactFile, []byte(reactContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test React file: %v", err)
	}
	
	// Set up scanner with parsers and registry
	parsers := []ComponentParser{
		NewVueParser(),
		NewReactParser(),
	}
	reg := registry.NewComponentMappingRegistry()
	scanner := NewComponentScanner(parsers, reg)
	
	t.Run("scan for forms finds q-form in Vue file", func(t *testing.T) {
		files := []string{vueFile}
		result, err := scanner.Scan(files, "form")
		
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		
		if result.TotalCount != 1 {
			t.Errorf("Expected 1 match, got %d", result.TotalCount)
		}
		
		if len(result.Matches) != 1 {
			t.Fatalf("Expected 1 match in results, got %d", len(result.Matches))
		}
		
		match := result.Matches[0]
		if match.ComponentName != "q-form" {
			t.Errorf("Expected component name 'q-form', got '%s'", match.ComponentName)
		}
		
		if match.ComponentType != "form" {
			t.Errorf("Expected component type 'form', got '%s'", match.ComponentType)
		}
		
		if result.ComponentType != "form" {
			t.Errorf("Expected result component type 'form', got '%s'", result.ComponentType)
		}
		
		if result.ScannedFiles != 1 {
			t.Errorf("Expected 1 scanned file, got %d", result.ScannedFiles)
		}
		
		if result.ScanTimeMs < 0 {
			t.Errorf("Expected positive scan time, got %d", result.ScanTimeMs)
		}
	})
	
	t.Run("scan for buttons finds Button in React file", func(t *testing.T) {
		files := []string{reactFile}
		result, err := scanner.Scan(files, "button")
		
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		
		if result.TotalCount != 1 {
			t.Errorf("Expected 1 match, got %d", result.TotalCount)
		}
		
		if len(result.Matches) != 1 {
			t.Fatalf("Expected 1 match in results, got %d", len(result.Matches))
		}
		
		match := result.Matches[0]
		if match.ComponentName != "Button" {
			t.Errorf("Expected component name 'Button', got '%s'", match.ComponentName)
		}
		
		if match.ComponentType != "button" {
			t.Errorf("Expected component type 'button', got '%s'", match.ComponentType)
		}
	})
	
	t.Run("scan multiple files concurrently", func(t *testing.T) {
		files := []string{vueFile, reactFile}
		result, err := scanner.Scan(files, "form")
		
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		
		// Should only find the form in Vue file, not the button in React file
		if result.TotalCount != 1 {
			t.Errorf("Expected 1 match, got %d", result.TotalCount)
		}
		
		if result.ScannedFiles != 2 {
			t.Errorf("Expected 2 scanned files, got %d", result.ScannedFiles)
		}
	})
	
	t.Run("scan with no matches returns empty result", func(t *testing.T) {
		files := []string{vueFile, reactFile}
		result, err := scanner.Scan(files, "dialog")
		
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		
		if result.TotalCount != 0 {
			t.Errorf("Expected 0 matches, got %d", result.TotalCount)
		}
		
		if len(result.Matches) != 0 {
			t.Errorf("Expected empty matches slice, got %d matches", len(result.Matches))
		}
		
		if result.ComponentType != "dialog" {
			t.Errorf("Expected component type 'dialog', got '%s'", result.ComponentType)
		}
	})
	
	t.Run("scan with non-existent file continues gracefully", func(t *testing.T) {
		nonExistentFile := filepath.Join(tempDir, "nonexistent.vue")
		files := []string{vueFile, nonExistentFile}
		result, err := scanner.Scan(files, "form")
		
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		
		// Should still find the form in the existing file
		if result.TotalCount != 1 {
			t.Errorf("Expected 1 match, got %d", result.TotalCount)
		}
	})
	
	t.Run("scan with unsupported file type skips file", func(t *testing.T) {
		unsupportedFile := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(unsupportedFile, []byte("some text"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		
		files := []string{unsupportedFile}
		result, err := scanner.Scan(files, "form")
		
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		
		if result.TotalCount != 0 {
			t.Errorf("Expected 0 matches for unsupported file, got %d", result.TotalCount)
		}
	})
}

func TestComponentScanner_filterByComponentType(t *testing.T) {
	reg := registry.NewComponentMappingRegistry()
	scanner := NewComponentScanner(nil, reg)
	
	t.Run("filters matches by component type", func(t *testing.T) {
		matches := []types.ComponentMatch{
			{ComponentName: "q-form", FilePath: "test.vue", Line: 1},
			{ComponentName: "q-btn", FilePath: "test.vue", Line: 2},
			{ComponentName: "Button", FilePath: "test.jsx", Line: 1},
		}
		
		filtered := scanner.filterByComponentType(matches, "form")
		
		if len(filtered) != 1 {
			t.Errorf("Expected 1 filtered match, got %d", len(filtered))
		}
		
		if filtered[0].ComponentName != "q-form" {
			t.Errorf("Expected 'q-form', got '%s'", filtered[0].ComponentName)
		}
		
		if filtered[0].ComponentType != "form" {
			t.Errorf("Expected component type 'form', got '%s'", filtered[0].ComponentType)
		}
	})
	
	t.Run("returns empty slice when no matches", func(t *testing.T) {
		matches := []types.ComponentMatch{
			{ComponentName: "q-btn", FilePath: "test.vue", Line: 1},
		}
		
		filtered := scanner.filterByComponentType(matches, "form")
		
		if len(filtered) != 0 {
			t.Errorf("Expected 0 filtered matches, got %d", len(filtered))
		}
	})
	
	t.Run("handles custom component types", func(t *testing.T) {
		matches := []types.ComponentMatch{
			{ComponentName: "CustomWidget", FilePath: "test.vue", Line: 1},
			{ComponentName: "OtherWidget", FilePath: "test.vue", Line: 2},
		}
		
		filtered := scanner.filterByComponentType(matches, "CustomWidget")
		
		if len(filtered) != 1 {
			t.Errorf("Expected 1 filtered match, got %d", len(filtered))
		}
		
		if filtered[0].ComponentName != "CustomWidget" {
			t.Errorf("Expected 'CustomWidget', got '%s'", filtered[0].ComponentName)
		}
	})
}
