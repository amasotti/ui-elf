package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"ui-elf/internal/types"
)

func TestShouldExcludeFile(t *testing.T) {
	service := NewFileDiscoveryService()

	tests := []struct {
		name     string
		filePath string
		filter   types.FileFilter
		expected bool
	}{
		{
			name:     "excludes node_modules",
			filePath: "src/node_modules/package/file.js",
			filter: types.FileFilter{
				ExcludePatterns: []string{"node_modules"},
			},
			expected: true,
		},
		{
			name:     "excludes test files",
			filePath: "src/components/Button.test.tsx",
			filter: types.FileFilter{
				ExcludePatterns: []string{".test.", ".spec."},
			},
			expected: true,
		},
		{
			name:     "excludes spec files",
			filePath: "src/components/Button.spec.tsx",
			filter: types.FileFilter{
				ExcludePatterns: []string{".test.", ".spec."},
			},
			expected: true,
		},
		{
			name:     "does not exclude regular files",
			filePath: "src/components/Button.tsx",
			filter: types.FileFilter{
				ExcludePatterns: []string{"node_modules", ".test.", ".spec."},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ShouldExcludeFile(tt.filePath, tt.filter)
			if result != tt.expected {
				t.Errorf("ShouldExcludeFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasValidExtension(t *testing.T) {
	service := NewFileDiscoveryService()

	tests := []struct {
		name       string
		filePath   string
		extensions []string
		expected   bool
	}{
		{
			name:       "matches .vue extension",
			filePath:   "src/components/Button.vue",
			extensions: []string{".vue", ".jsx", ".tsx"},
			expected:   true,
		},
		{
			name:       "matches .jsx extension",
			filePath:   "src/components/Button.jsx",
			extensions: []string{".vue", ".jsx", ".tsx"},
			expected:   true,
		},
		{
			name:       "matches .tsx extension",
			filePath:   "src/components/Button.tsx",
			extensions: []string{".vue", ".jsx", ".tsx"},
			expected:   true,
		},
		{
			name:       "does not match .js extension",
			filePath:   "src/components/Button.js",
			extensions: []string{".vue", ".jsx", ".tsx"},
			expected:   false,
		},
		{
			name:       "empty extensions list matches all",
			filePath:   "src/components/Button.js",
			extensions: []string{},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.hasValidExtension(tt.filePath, tt.extensions)
			if result != tt.expected {
				t.Errorf("hasValidExtension() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDiscoverFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"src/components/Button.vue",
		"src/components/Form.jsx",
		"src/components/Dialog.tsx",
		"src/components/Button.test.tsx",
		"node_modules/package/file.js",
		"tests/unit/test.spec.js",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(fullPath)

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	service := NewFileDiscoveryService()

	t.Run("discovers files with extension filter", func(t *testing.T) {
		filter := types.FileFilter{
			ExcludePatterns: []string{"node_modules", ".test.", ".spec."},
			FileExtensions:  []string{".vue", ".jsx", ".tsx"},
		}

		files, err := service.DiscoverFiles(tmpDir, filter)
		if err != nil {
			t.Fatalf("DiscoverFiles() error = %v", err)
		}

		// Should find 3 files: Button.vue, Form.jsx, Dialog.tsx
		if len(files) != 3 {
			t.Errorf("DiscoverFiles() found %d files, want 3", len(files))
		}
	})

	t.Run("respects include directories filter", func(t *testing.T) {
		filter := types.FileFilter{
			ExcludePatterns:    []string{"node_modules", ".test.", ".spec."},
			IncludeDirectories: []string{"src/components"},
			FileExtensions:     []string{".vue", ".jsx", ".tsx"},
		}

		files, err := service.DiscoverFiles(tmpDir, filter)
		if err != nil {
			t.Fatalf("DiscoverFiles() error = %v", err)
		}

		// Should find 3 files in src/components
		if len(files) != 3 {
			t.Errorf("DiscoverFiles() found %d files, want 3", len(files))
		}

		// Verify all files are in src/components
		for _, file := range files {
			relPath, _ := filepath.Rel(tmpDir, file)
			if !filepath.HasPrefix(relPath, "src/components") {
				t.Errorf("File %s is not in src/components", relPath)
			}
		}
	})
}
