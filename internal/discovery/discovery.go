// Package discovery handles file system traversals to find relevant files.
package discovery

import (
	"os"
	"path/filepath"
	"strings"

	"ui-elf/internal/types"
)

// FileDiscoveryService handles file discovery with filtering
type FileDiscoveryService struct{}

// NewFileDiscoveryService creates a new FileDiscoveryService
func NewFileDiscoveryService() *FileDiscoveryService {
	return &FileDiscoveryService{}
}

// DiscoverFiles traverses the directory tree and returns files matching the filter criteria
func (s *FileDiscoveryService) DiscoverFiles(rootDir string, filter types.FileFilter) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file should be excluded
		if s.ShouldExcludeFile(path, filter) {
			return nil
		}

		// Check if file has a valid extension
		if !s.hasValidExtension(path, filter.FileExtensions) {
			return nil
		}

		// If include directories are specified, check if file is in one of them
		if len(filter.IncludeDirectories) > 0 {
			if !s.isInIncludedDirectory(path, rootDir, filter.IncludeDirectories) {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

// ShouldExcludeFile checks if a file should be excluded based on filter patterns
func (s *FileDiscoveryService) ShouldExcludeFile(filePath string, filter types.FileFilter) bool {
	for _, pattern := range filter.ExcludePatterns {
		if s.matchesPattern(filePath, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a file path matches an exclusion pattern
func (s *FileDiscoveryService) matchesPattern(filePath string, pattern string) bool {
	// Normalize path separators
	normalizedPath := filepath.ToSlash(filePath)

	// Check if path contains the pattern
	if strings.Contains(normalizedPath, pattern) {
		return true
	}

	// Check if any directory component matches the pattern
	parts := strings.Split(normalizedPath, "/")
	for _, part := range parts {
		if part == pattern {
			return true
		}
	}

	return false
}

// hasValidExtension checks if a file has one of the valid extensions
func (s *FileDiscoveryService) hasValidExtension(filePath string, extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}

	ext := filepath.Ext(filePath)
	for _, validExt := range extensions {
		if ext == validExt {
			return true
		}
	}

	return false
}

// isInIncludedDirectory checks if a file is within one of the included directories
func (s *FileDiscoveryService) isInIncludedDirectory(filePath string, rootDir string, includeDirectories []string) bool {
	// Get relative path from root
	relPath, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		return false
	}

	normalizedRelPath := filepath.ToSlash(relPath)

	for _, includeDir := range includeDirectories {
		normalizedIncludeDir := filepath.ToSlash(includeDir)

		// Check if file is in the included directory or its subdirectories
		if strings.HasPrefix(normalizedRelPath, normalizedIncludeDir+"/") ||
			normalizedRelPath == normalizedIncludeDir {
			return true
		}
	}

	return false
}
