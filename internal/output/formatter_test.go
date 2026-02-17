package output

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"ui-elf/internal/types"
)

func TestFormatTerminal(t *testing.T) {
	formatter := NewOutputFormatter()
	
	t.Run("formats result with matches", func(t *testing.T) {
		result := &types.ScanResult{
			Matches: []types.ComponentMatch{
				{
					FilePath:      "src/components/Form.vue",
					Line:          10,
					ComponentName: "q-form",
					ComponentType: "form",
				},
				{
					FilePath:      "src/pages/Login.vue",
					Line:          25,
					ComponentName: "form",
					ComponentType: "form",
				},
			},
			TotalCount:    2,
			ScanTimeMs:    150,
			ComponentType: "form",
			ScannedFiles:  50,
		}
		
		output := formatter.FormatTerminal(result)
		
		// Verify output contains key information
		if !strings.Contains(output, "Component Finder Results - form") {
			t.Error("Output should contain component type header")
		}
		if !strings.Contains(output, "src/components/Form.vue") {
			t.Error("Output should contain first file path")
		}
		if !strings.Contains(output, "src/pages/Login.vue") {
			t.Error("Output should contain second file path")
		}
		if !strings.Contains(output, "line 10") {
			t.Error("Output should contain line number")
		}
		if !strings.Contains(output, "q-form") {
			t.Error("Output should contain component name")
		}
		if !strings.Contains(output, "Total components found: 2") {
			t.Error("Output should contain total count")
		}
		if !strings.Contains(output, "Files scanned: 50") {
			t.Error("Output should contain scanned files count")
		}
		if !strings.Contains(output, "Scan time: 150ms") {
			t.Error("Output should contain scan time")
		}
	})
	
	t.Run("formats result with no matches", func(t *testing.T) {
		result := &types.ScanResult{
			Matches:       []types.ComponentMatch{},
			TotalCount:    0,
			ScanTimeMs:    100,
			ComponentType: "button",
			ScannedFiles:  30,
		}
		
		output := formatter.FormatTerminal(result)
		
		if !strings.Contains(output, "No components found") {
			t.Error("Output should indicate no components found")
		}
		if !strings.Contains(output, "Total components found: 0") {
			t.Error("Output should show zero count")
		}
		if !strings.Contains(output, "Scan time: 100ms") {
			t.Error("Output should contain scan time")
		}
	})
}

func TestFormatJSON(t *testing.T) {
	formatter := NewOutputFormatter()
	
	t.Run("formats result as valid JSON", func(t *testing.T) {
		result := &types.ScanResult{
			Matches: []types.ComponentMatch{
				{
					FilePath:      "src/App.tsx",
					Line:          15,
					ComponentName: "Button",
					ComponentType: "button",
				},
			},
			TotalCount:    1,
			ScanTimeMs:    200,
			ComponentType: "button",
			ScannedFiles:  20,
		}
		
		jsonStr, err := formatter.FormatJSON(result)
		if err != nil {
			t.Fatalf("FormatJSON failed: %v", err)
		}
		
		// Verify it's valid JSON by unmarshaling
		var parsed types.ScanResult
		if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
			t.Fatalf("Generated JSON is invalid: %v", err)
		}
		
		// Verify data integrity
		if parsed.TotalCount != 1 {
			t.Errorf("Expected TotalCount 1, got %d", parsed.TotalCount)
		}
		if parsed.ComponentType != "button" {
			t.Errorf("Expected ComponentType 'button', got %s", parsed.ComponentType)
		}
		if parsed.ScanTimeMs != 200 {
			t.Errorf("Expected ScanTimeMs 200, got %d", parsed.ScanTimeMs)
		}
		if parsed.ScannedFiles != 20 {
			t.Errorf("Expected ScannedFiles 20, got %d", parsed.ScannedFiles)
		}
		if len(parsed.Matches) != 1 {
			t.Fatalf("Expected 1 match, got %d", len(parsed.Matches))
		}
		if parsed.Matches[0].FilePath != "src/App.tsx" {
			t.Errorf("Expected FilePath 'src/App.tsx', got %s", parsed.Matches[0].FilePath)
		}
	})
}

func TestWrite(t *testing.T) {
	formatter := NewOutputFormatter()
	result := &types.ScanResult{
		Matches: []types.ComponentMatch{
			{
				FilePath:      "test.vue",
				Line:          5,
				ComponentName: "q-dialog",
				ComponentType: "dialog",
			},
		},
		TotalCount:    1,
		ScanTimeMs:    50,
		ComponentType: "dialog",
		ScannedFiles:  10,
	}
	
	t.Run("writes JSON to file", func(t *testing.T) {
		tmpFile := "test-output.json"
		defer os.Remove(tmpFile)
		
		err := formatter.Write(result, "json", tmpFile)
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		
		// Verify file was created
		if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
			t.Error("JSON file was not created")
		}
		
		// Verify file content
		content, err := os.ReadFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}
		
		var parsed types.ScanResult
		if err := json.Unmarshal(content, &parsed); err != nil {
			t.Fatalf("Output file contains invalid JSON: %v", err)
		}
		
		if parsed.TotalCount != 1 {
			t.Errorf("Expected TotalCount 1, got %d", parsed.TotalCount)
		}
	})
	
	t.Run("uses default filename when not specified", func(t *testing.T) {
		defaultFile := "component-finder-results.json"
		defer os.Remove(defaultFile)
		
		err := formatter.Write(result, "json", "")
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		
		// Verify default file was created
		if _, err := os.Stat(defaultFile); os.IsNotExist(err) {
			t.Error("Default JSON file was not created")
		}
	})
	
	t.Run("returns error for unsupported format", func(t *testing.T) {
		err := formatter.Write(result, "invalid", "")
		if err == nil {
			t.Error("Expected error for unsupported format")
		}
		if !strings.Contains(err.Error(), "unsupported output format") {
			t.Errorf("Expected 'unsupported output format' error, got: %v", err)
		}
	})
}
